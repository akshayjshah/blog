# Language Use on GitHub

Do CoffeeScript aficionados write less vanilla JavaScript? Do systems hackers
using Go also do front-end work in ActionScript? Most programmers have some
intuitions about these questions---*but what does the data say*? Spurred on by
a six-month-old Twitter [conversation][twitter-convo], I decided to find out.
Using data from 2012, I charted the relationships between the 25 most popular
languages on GitHub: <figure><img
src="/img/language-use-on-github/spearman_language_correlation.svg"
alt="Language Correlation on GitHub"></figure>

Each square in the chart measures the [rank correlation][wiki-rank-correlation]
between two languages; positive correlations are in blue, and negative
correlations are in red. A little more than halfway down the first column is a
medium-red square showing the correlation between Go and ActionScript. Because
the square is red, we know that as users write more Go, they write less
ActionScript. From the intensity of the red, we know that this rule of thumb is
fairly reliable. Rank correlations are symmetric, so ActionScript fans also
write less Go (and there's an identically-colored square in the bottom row to
show it). Take a moment to find your favorite languages, then read on for the
details!

## Data Collection

GitHub's own [Brian Doll][brian-doll] published a [data set][brian-doll-data]
titled "Programming Language Correlations," but as an astute commenter pointed
out, it's really a set of conditional probabilities. For our question,
that's an important distinction---while conditional probabilities let us say,
"87.9% of CoffeeScript programmers also code in Ruby," they *don't* allow us to
say, "People who write more CoffeeScript also tend to write more Ruby." To
tackle our question, we'll need access to some more granular data.

Rather than banging on the API, we can use the [GitHub
Archive][github-archive]. This fantastic resource archives every public GitHub
event and makes the whole data set accessible via [Google BigQuery][big-query].
BigQuery has a web-based console and a comfortably SQL-like query language, so
it's easy to get the data we need (all the code in this post is also in a
single [Gist][gist]):

```sql
select actor, repository_language, count(repository_language) as pushes
from [githubarchive:github.timeline]
where type='PushEvent'
    and repository_language != ''
    and PARSE_UTC_USEC(created_at) >= PARSE_UTC_USEC('2012-01-01 00:00:00')
    and PARSE_UTC_USEC(created_at) < PARSE_UTC_USEC('2013-01-01 00:00:00')
group by actor, repository_language;
```

The results of this query are in a stacked format, where each combination of
user and language is on a separate row:

<table>
  <thead>
    <tr>
      <th>actor</th>
      <th>repository_language</th>
      <th>pushes</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Joe</td>
      <td>CoffeeScript</td>
      <td>24</td>
    </tr>
    <tr>
      <td>Joe</td>
      <td>Ruby</td>
      <td>72</td>
    </tr>
    <tr>
      <td>Sue</td>
      <td>CoffeeScript</td>
      <td>7</td>
    </tr>
    <tr>
      <td>Joe</td>
      <td>Go</td>
      <td>1</td>
    </tr>
    <tr>
      <td>Sue</td>
      <td>Ruby</td>
      <td>48</td>
    </tr>
  </tbody>
</table>

Stacked formats are often convenient in database schemas, but they're not very
useful for analysis. We'd rather unstack the data so that there's one row per
user and one column per language:

<table>
  <thead>
    <tr>
      <th>actor</th>
      <th>CoffeeScript</th>
      <th>Go</th>
      <th>Ruby</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Joe</td>
      <td>24</td>
      <td>1</td>
      <td>72</td>
    </tr>
    <tr>
      <td>Sue</td>
      <td>0</td>
      <td>0</td>
      <td>48</td>
    </tr>
  </tbody>
</table>

With the help of the brilliant [pandas][] library, we can perform this
transformation with one Python command:

```python
import pandas as pd
pushes = pd.read_csv('stacked_language_by_user.csv').pivot(
    index='actor',
    columns='repository_language',
    values='pushes')
```

Exporting these results from BigQuery was an enormous pain (and required a paid
account), so I'll keep a zipped copy of the [stacked][stacked-csv] and
[unstacked][unstacked-csv] results available.

## Analysis

GitHub recognizes *lots* of different languages, including some that are fairly
obscure, so our unstacked data set has too many columns to visualize. (As
[conjugateprior](http://conjugateprior.org) noted in the now-defunct comments,
this calculation ignores any rows with missing data: for example, the
Python-Ruby correlation ignores any users who haven't used both Python and
Ruby. We could fill the missing values with zeroes, or we could also calculate
significance for each correlation. Corey has already done the former---check out his
[updated code](https://gist.github.com/coyotebush/5379476) and
[plot](http://coreyford.name/2013/04/13/github-language-correlations.html).)
Let's just keep the most popular languages:

```python
import numpy as np
popular = pushes.select(lambda x: np.sum(pushes[x]) > 50000, axis=1)
```

Now that our data's formatted and filtered, it's time to actually calculate our
correlation matrix and draw a plot. Again, pandas makes the number-crunching
ridiculously simple:

```python
import matplotlib.pyplot as plt

def plot_correlation(dataframe, filename, title='', corr_type=''):
    lang_names = dataframe.columns.tolist()
    tick_indices = np.arange(0.5, len(lang_names) + 0.5)
    plt.figure()
    plt.pcolor(dataframe.values, cmap='RdBu', vmin=-1, vmax=1)
    colorbar = plt.colorbar()
    colorbar.set_label(corr_type)
    plt.title(title)
    plt.xticks(tick_indices, lang_names, rotation='vertical')
    plt.yticks(tick_indices, lang_names)
    plt.savefig(filename)

spearman_corr = popular.corr(method='spearman')
plot_correlation(
    spearman_corr,
    'spearman_language_correlation.svg',
    title='2012 GitHub Language Correlations',
    corr_type='Spearman\'s Rank Correlation')
```

It's better to use [Spearman's rank correlation][wiki-rank-correlation] here
instead of the usual [Pearson correlation][wiki-pearson-correlation] for two
reasons:

* We don't really care whether the relationship between languages is strictly
  linear.
* There are quite a few outliers in our data set, and rank correlations are
  less distorted by these outliers.

If that doesn't convince you, it's easy to calculate the Pearson correlation---it's
the default in pandas, so removing the ``method='spearman'`` above should do
the trick. If you're impatient, you can just peek at the
[results][pearson-plot].

<h2 id="conclusions">Conclusions</h2>
<img src="/img/language-use-on-github/spearman_language_correlation.svg"
alt="Language Correlation on GitHub">

The most striking thing about this chart is its *blueness*. Despite our
tribalism, writing scads of C# doesn't make programmers any less likely to hack
on some R. Even PHP, perhaps the [most hated][php] programming language on
earth, has a slight positive correlation with Haskell. After seeing so many
flamewars in forums, on mailing lists, and even in person, I expected language
communities to be more insular. I'm particularly surprised by the positive
correlations between the languages associated with proprietary platforms (C#,
Objective-C, and ActionScript) and the traditionally open-source languages.

Not surprisingly, special-purpose languages are the exception to this rule. R,
Matlab, and Puppet have more strong correlations (both positive and negative)
than the norm, likely because of their niche roles in data analysis and devops.

Like any analysis project, this one comes with a few caveats:

* GitHub pushes aren't a perfect measure of activity. Then again, neither are
  commits, lines of code changed, or anything else I've heard of.
* This data only considers public projects on GitHub, many of which are open
  source. Open-source programmers, and projects, may behave quite differently
  from their closed-source counterparts.
* Correlation isn't causation.

I've only scratched the surface here---if you've got some ideas, [let me
know](mailto:akshay@akshayshah.org)!

[big-query]: https://developers.google.com/bigquery/
[brian-doll-data]: https://gist.github.com/briandoll/e0637fff9c8eec988528
[brian-doll]: https://github.com/briandoll
[gist]: https://gist.github.com/akshayjshah/4772174
[github-archive]: http://www.githubarchive.org/
[pandas]: http://pandas.pydata.org/
[pearson-plot]: /img/language-use-on-github/pearson_language_correlation.svg
[php]: http://me.veekun.com/blog/2012/04/09/php-a-fractal-of-bad-design/
[stacked-csv]: /img/language-use-on-github/stacked_language_by_user.zip
[twitter-convo]: https://twitter.com/misc/status/235167513833525249
[unstacked-csv]: /img/language-use-on-github/unstacked_language_by_user.zip
[wiki-pearson-correlation]: http://en.wikipedia.org/wiki/Pearson_product-moment_correlation_coefficient
[wiki-rank-correlation]: http://en.wikipedia.org/wiki/Spearman%27s_rank_correlation_coefficient
