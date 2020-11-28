{{ range .Result }}
<div class="card">
    <div class="card-header bg-secondary" style="height:100px;">
        <div class="container">
            <div class="row">
                
                <div class="col text-truncate">
                    <h2 class="font-weight-normal text-wrap text-white text-capitalize" title="{{.Title}}"><small>{{.Title}}</small></h2>
                </div>
                
                <div class="col-3 text-right">
                    <a href="{{.ExternalURL}}" target="_blank">
                        <img src="/img/{{ .Provider }}-100x30.png" class="image-fluid" alt="{{ .Provider }} logo small" />
                    </a>
                </div>

            </div>
        </div>
    </div>
  <div class="card-body">
    <div class="album py-5 bg-white w-100 mx-auto">
        <div class="container embed-responsive embed-responsive-16by9">
            {{ .Embed }}
        </div>
    </div>
    {{ range .Pornstars }}
    <a href="/{{.}}.html" class="badge badge-dark" alt="{{.}}">{{.}}</a>
    {{ end }}
    {{ range .Categories }}
    <a href="/{{.}}.html" class="badge badge-light" alt="{{.}}">{{.}}</a>
    {{ end }}
    {{ range .Tags }}
    <a href="/{{.}}.html" class="badge badge-light" alt="{{.}}">{{.}}</a>
    {{ end }}
  </div>
    <div class="card-footer bg-dark text-white">
        <small>{{ .PublishDate }}</small>
        <small class="float-right">{{ .Duration }}</small>
    </div>
</div>
{{ end }}