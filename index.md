# Akshay Shah

<img alt="" src="/img/akshay-300x300.png" height="300" width="300">

Hi there! I'm Akshay, and you've stumbled across my oft-neglected
blog.  I've been a medical student, a public school teacher, an EMT,
a spammer, and a programmer; these days, I'm the cofounder of my own little
startup. Shoot me an [email](mailto:akshay@akshayshah.org) and say hello,
or [read a bit more about me](/colophon/).

# Writing

<ol class="post-list">
{{ range . }}
<li>
  <span class="post-date">{{ .Published }} &raquo;</span>
  <a href="{{ .Link }}">{{ .Title }}</a>
</li>
{{ end }}
</ol>
