# Akshay Shah

<img alt="" src="/static/headshot-2023-300x300.webp" height="300px" width="300px">

Hi there! I'm Akshay, and you've stumbled across my oft-neglected blog. I've
been an infrastructure and platform engineer for 15 years; before that, I was a
medical student, a public school teacher, an EMT, a fax spammer, and a data
scientist. Shoot me an [email](mailto:akshay@akshayshah.org) and say hello, or
[read a bit more about me](/colophon/).

# Writing

<ol class="post-list">
{{ range .Posts }}
<li>
  <span class="post-date">{{ .Published }} &nbsp;&nbsp;</span>
  <a href="{{ .Link }}">{{ .Title }}</a>{{ if .Via }} ({{ .Via}}){{ end }}
</li>
{{ end }}
</ol>

# Cooking

<ol class="post-list">
{{ range .Recipes }}
<li>
  <span class="post-date">{{ .Published }} &nbsp;&nbsp;</span>
  <a href="{{ .Link }}">{{ .Title }}</a>
</li>
{{ end }}
</ol>
