{{ range .Result }}
<div class="card">
    <div class="card-header bg-secondary" style="height:100px;">
        <div class="container">
            <div class="row">
                
                <div class="col text-truncate">
                    <a href="/{{ .Provider }}/{{ .ID }}.html" class="text-decoration-none">
                        <h2 class="font-weight-normal text-wrap text-white text-capitalize" title="{{.Title}}"><small>{{.Title}}</small></h2>
                    </a>  
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
    {{ .Embed }}
    <h5 class="card-title">Special title treatment</h5>
    <p class="card-text">With supporting text below as a natural lead-in to additional content.</p>
    <a href="#" class="btn btn-primary">Go somewhere</a>
  </div>
</div>
{{ end }}