---
title: "Akshay Shah"
description: "After building infrastructure at Uber, Microsoft, and a pile of startups, I'm now the Field CTO at Antithesis. I typically write about developer tools and infrastructure, with an occasionally foray into cooking."
created: {{ .Created }}
updated: {{ .Updated }}
hidden: true
hide_home: true
hide_dates: true
priority: 0.8
---

<img src="/static/headshot-2025-450x450.webp" alt="" class="profile-pic" height="300px" width="300px" fetchpriority="high">

Hi there! I'm Akshay, and you've stumbled across my oft-neglected blog. After
building infrastructure at Uber, Microsoft, and a pile of startups, I'm now the
Field CTO at [Antithesis](https://antithesis.com). We make the world's best
tools for testing distributed systems.

Shoot me an [email](mailto:akshay@akshayshah.org) and say hello, or
[read a bit more about me](/colophon/).

## Writing &amp; speaking

{{ range .Posts }}
<div class="post-row">
  <span class="post-date">{{ .Created }}</span>
  <a href="{{ .Link }}">{{ .Title }}</a>{{ if .Via }} ({{ .Via}}){{ end }}
</div>
{{ end }}

## Cooking

{{ range .Recipes }}
<div class="post-row">
  <span class="post-date">{{ .Created }}</span>
  <a href="{{ .Link }}">{{ .Title }}</a>
</div>
{{ end }}
