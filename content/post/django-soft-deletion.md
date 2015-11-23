+++
title = "Soft Deletion in Django"
date = "2013-05-09"
+++

Like any self-respecting data nerd, I find deleting database records abhorrent.
What happens if I need to resurrect those records later? What if I want to run
a survival analysis? Django's ORM doesn't offer any out-of-the-box support for
soft-deletion, but it's not difficult to preserve your data by overriding a few
methods and flipping a Boolean field instead of actually deleting anything.
After living with that system for a while, though, I'm convinced that it's
inadequate; in this post, I'll explain its principal shortcoming and propose a
slightly more complex, but vastly better, alternative.
{{< update "django-livefield" >}}
I&rsquo;ve open-sourced this soft-deletion
plug-in as <code>django-livefield</code>. You can install it from <a
href="https://pypi.python.org/pypi/django-livefield/" title="django-livefield
on PyPI">PyPI</a> or check out the code on <a
href="https://github.com/hearsaycorp/django-livefield" title="django-livefield
on GitHub">GitHub</a>. Try it out and let me know what you think!
{{< /update >}}


## Naive Soft-Deletion

At first blush, soft-deletion seems embarrassingly simple. If you're like me,
a system like this probably springs to mind:

```python
from django.db import models
from django.db.models.query import QuerySet


class SoftDeletionQuerySet(QuerySet):
    def delete(self):
        # Bulk delete bypasses individual objects' delete methods.
        return super(SoftDeletionQuerySet, self).update(alive=False)

    def hard_delete(self):
        return super(SoftDeletionQuerySet, self).delete()

    def alive(self):
        return self.filter(alive=True)

    def dead(self):
        return self.exclude(alive=True)


class SoftDeletionManager(models.Manager):
    def __init__(self, *args, **kwargs):
        self.alive_only = kwargs.pop('alive_only', True)
        super(SoftDeletionManager, self).__init__(*args, **kwargs)

    def get_queryset(self):
        if self.alive_only:
            return SoftDeletionQuerySet(self.model).filter(alive=True)
        return SoftDeletionQuerySet(self.model)

    def hard_delete(self):
        return self.get_queryset().hard_delete()


class SoftDeletionModel(models.Model):
    alive = models.BooleanField(default=True)

    objects = SoftDeletionManager()
    all_objects = SoftDeletionManager(alive_only=False)

    class Meta:
        abstract = True

    def delete(self):
        self.alive = False
        self.save()

    def hard_delete(self):
        super(SoftDeletionModel, self).delete()
```

This approach is straightforward and readable, and for nearly two years it
worked well for us at [Hearsay Social](http://hearsaysocial.com/careers/).
*However, it inevitably leads to data corruption.*

The problem is simple: using a Boolean to store deletion status makes it
impossible to enforce uniqueness constraints in your database. Let's say you're
storing user records which should have unique email addresses; with this
soft-deletion scheme, you can only have one active record for
"betty@smith.com". Including deletion status in your constraint lets you keep
both a soft-deleted and an active record with the same email address, but then
you're out of luck -- any attempt to delete another record for Betty will throw
an ``IntegrityError``. Luckily, there's a better way.

## The Null Solution

At the database level, there's a straightforward solution to this problem
(though I didn't learn about it until a few months ago): store soft-deleted
records with nulls in the ``alive`` column. As mandated by the ANSI SQL
standard, MySQL, Postgres, and SQLite treat each null as a unique snowflake.

However, creating a Django field with this behavior is a little tricky because
we want to forbid ``False`` values in the database (allowing only ``True`` and
``NULL``).  Here's what I came up with:

```python
from django.db import models

class LiveField(models.Field):
    '''Similar to a BooleanField, but stores False as NULL.

    '''
    description = 'Soft-deletion status'
    __metaclass__ = models.SubfieldBase

    def __init__(self):
        super(LiveField, self).__init__(default=True, null=True)

    def get_internal_type(self):
        # Create DB column as though for a NullBooleanField.
        return 'NullBooleanField'

    def get_prep_value(self, value):
        # Convert in-Python value to value we'll store in DB
        if value:
            return 1
        return None

    def to_python(self, value):
        # Misleading name, since type coercion also occurs when
        # assigning a value to the field in Python.
        return bool(value)

    def get_prep_lookup(self, lookup_type, value):
        # Filters with .alive=False won't work, so
        # raise a helpful exception instead.
        if lookup_type == 'exact' and not value:
            msg = ("%(model)s doesn't support filters with "
                "%(field)s=False. Use a filter with "
                "%(field)s=None or an exclude with "
                "%(field)s=True instead.")
            raise TypeError(msg % {
                'model': self.model.__name__,
                'field': self.name})

        return super(LiveField, self).get_prep_lookup(lookup_type, value)
```

This is a drop-in replacement for Django's stock ``BooleanField`` in the
abstract model above, but under the covers it stores falsy values as nulls. At
Hearsay, we just finished migrating all our models to use ``LiveField``, and
it's been a huge help already -- having the option to simultaneously support
soft-deletion and uniqueness constraints keeps our application code *and* our
data clean.

Curious how to test your shiny new soft-deletion field? Check out my post on
[testing custom Django fields](/post/testing-django-fields) for some tips, or
check out the actual test setup on
[GitHub](https://github.com/hearsaycorp/django-livefield).
