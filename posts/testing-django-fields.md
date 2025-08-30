# Testing Django Fields

I love the flexibility that custom Django fields, abstract models, managers,
and querysets offer, but unit testing them is a pain. Ideally, the tests for
custom Django fields should be completely isolated from the models that use the
fields in production; deciding that, for example, my ``User`` model no longer
needs to support soft deletion shouldn't affect the tests for the [soft-deletion
field][] itself.

The [most common approach][so] to this problem is simple, if annoying: declare
all test-specific models in ``test/models.py``, but don't include the test app
in ``INSTALLED APPS``. In your test suite's setup method, monkey-patch your
settings to include the test app and run Django's ``syncdb`` command, and
un-patch your settings in the teardown method. Dynamically altering your
settings in the test suite keeps your production database clean --- a rogue
``syncdb`` won't suddenly create dozens of useless new tables. My biggest gripe
with this approach, though, is that it forces you to separate your test code
into two files. The tests become much harder to read, and the file of test
models inevitably becomes a crufty mess.

After a few months of low-level frustration, I finally came up with a better
solution. By making all my test models inherit from this abstract model, I can
have it all: no raw SQL, no test tables in production, and model definitions
alongside my test code.

```python
from django.core.management.color import no_style
from django.db import connection, models


class TestModel(models.Model):

    class Meta:
        abstract = True

    @classmethod
    def create_table(cls):
        # Cribbed from Django's management commands.
        raw_sql, refs = connection.creation.sql_create_model(
            cls,
            no_style(),
            [])
        create_sql = u'\n'.join(raw_sql).encode('utf-8')
        cls.delete_table()
        cursor = connection.cursor()
        try:
            cursor.execute(create_sql)
        finally:
            cursor.close()

    @classmethod
    def delete_table(cls):
        cursor = connection.cursor()
        try:
            cursor.execute('DROP TABLE IF EXISTS %s' % cls._meta.db_table)
        except:
            # Catch anything backend-specific here.
            # (E.g., MySQLdb raises a warning if the table didn't exist.)
            pass
        finally:
            cursor.close()
```

To avoid boilerplate table management in my test setup and teardown code, I
added a little functionality to Django's built-in ``TestCase``.

```python
from django.test import TestCase


class ModelTestCase(TestCase):
    temporary_models = tuple()

    def setUp(self):
        self._map_over_temporary_models('create_table')
        super(ModelTestCase, self).setUp()

    def tearDown(self):
        self._map_over_temporary_models('delete_table')
        super(ModelTestCase, self).tearDown()

    def _map_over_temporary_models(self, method_name):
        for m in self.temporary_models:
            try:
                getattr(m, method_name)()
            except AttributeError:
                raise TypeError("%s doesn't support table mgmt." % m)
```

Looking for an example? Here's a section of the test suite for my
[soft-deletion field][]:

```python
from django.db import IntegrityError, models

from myproject.soft_deletion.models import SoftDeletionModel
from myproject.test.models import TestModel
from myproject.test.testcase import ModelTestCase


class Person(SoftDeletionModel, TestModel):
    name = models.CharField(max_length=20)

    class Meta:
        unique_together = ('name', 'alive')


class SoftDeletionTests(ModelTestCase):
    temporary_models = (Person,)

    def test_inits_alive(self):
        p = Person.objects.create(name='Alive')
        self.assertTrue(p.alive)

    def test_allows_many_deleted_with_same_name(self):
        Person.objects.create(name='Akshay').delete()
        Person.objects.create(name='Akshay').delete()

        # One un-deleted Akshay is okay.
        Person.objects.create(name='Akshay')
        self.assertEqual(Person.all_objects.count(), 3)

        # Resurrecting one of the dupes violates constraint.
        first = Person.all_objects.all()[0]
        first.alive = True
        self.assertRaises(IntegrityError, first.save)
```

Questions? Have a better idea? [Let me know](mailto:akshay@akshayshah.org)!

[so]: http://stackoverflow.com/questions/502916/django-how-to-create-a-model-dynamically-just-for-testing
[soft-deletion field]: /soft-deletion-in-django
