    <head>
        <meta charset="utf-8" />
        <title>{{.PageTitle}}</title>
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <meta name="description" content="{{.PageMetaDesc}}">
        {{ if .Result -}}
        <meta name="keywords" content="{{ range (index .Result 0).Tags }}{{ . }},{{ end }}just-tit">
        {{ end -}}
        <meta name="author" content="@dsmatilla" />
        <meta name="theme-color" content="#FFFFFF"/>
        <link rel="apple-touch-icon" href="/img/icon-192x192.png"/>
        <link rel="manifest" href="/manifest.json"/>
        <link rel="shortcut icon" type="image/png" href="/img/favicon.png"/>
        <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/css/bootstrap.min.css" integrity="sha384-TX8t27EcRE3e/ihU7zmQxVncDAy5uIKz4rEkgIXeMed4M0jlfIDPvg6uqKI2xXr2" crossorigin="anonymous">
        {{ if .Result -}}
        {{ if eq (index .Result 0).Type  "single" -}}
        <meta property="og:type" content="video">
        <meta property="og:title" content="{{.PageTitle}}"/>
        <meta property="og:image" content="{{(index .Result 0).Domain}}{{ToImageProxy (index .Result 0).Thumb}}"/>
        <meta property="og:url" content="{{(index .Result 0).URL}}"/>
        <meta name="twitter:card" content="player"/>
        <meta name="twitter:site" content="@Just_Tit"/>
        <meta name="twitter:title" content="{{.PageTitle}}"/>
        <meta name="twitter:image" content="{{(index .Result 0).Domain}}{{ToImageProxy (index .Result 0).Thumb}}"/>
        <meta name="twitter:player" content="{{(index .Result 0).URL}}?tp=true"/>
        {{ end -}}
        {{ end -}}
    </head>