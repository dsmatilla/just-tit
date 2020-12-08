{{- range .SearchResult }}
                        <div class="col-md-4 d-flex">
                            <div class="card border-secondary mb-3">
                                <div class="card-header bg-secondary" style="height:90px;">
                                    <div class="container">
                                        <div class="row">
                                            
                                            <div class="col text-truncate">
                                                <a href="/{{ .Provider }}/{{ .ID }}.html" class="text-decoration-none">
                                                    <h6 class="font-weight-normal text-wrap text-white text-capitalize" title="{{.Title}}"><small>{{.Title}}</small></h6>
                                                </a>  
                                            </div>
                                            
                                            <div class="col-4 text-right">
                                                <a href="{{.ExternalURL}}" target="_blank">
                                                    <img src="/img/{{ .Provider }}-100x30.png" class="image-fluid" alt="{{ .Provider }} logo small" loading="lazy" alt="{{ .Provider }} logo"/>
                                                </a>
                                            </div>

                                        </div>
                                    </div>
                                </div>
                                <div class="card-body text-secondary bg-light">
                                    <div id="carousel{{ .ID }}" class="carousel slide" data-ride="carousel" data-interval="false">
                                        <div class="carousel-inner">
                                            <a href="/{{ .Provider }}/{{ .ID }}.html" class="stretched-link"></a>
                                            <div class="carousel-item active">
                                                <img class="d-block w-100 h-100" src="{{ToImageProxy .Thumb}}" alt="{{ .Title }}" loading="lazy">
                                            </div>
                                            {{ $alt := .Title }}{{ range .Thumbs -}}
                                            <div class="carousel-item">
                                                <img class="d-block w-100 h-100" src="{{ToImageProxy .}}" alt="{{$alt}}" loading="lazy">
                                            </div>
                                            {{ end -}}
                                        </div>
                                        <a class="carousel-control-prev" href="#carousel{{ .ID }}" role="button" data-slide="prev">
                                            <span class="carousel-control-prev-icon" aria-hidden="true"></span>
                                            <span class="sr-only">Previous</span>
                                        </a>
                                        <a class="carousel-control-next" href="#carousel{{ .ID }}" role="button" data-slide="next">
                                            <span class="carousel-control-next-icon" aria-hidden="true"></span>
                                            <span class="sr-only">Next</span>
                                        </a>
                                    </div>
                                    {{ range .Pornstars -}}
                                    <a href="/{{.}}.html" class="badge badge-dark" alt="{{.}}">{{.}}</a>
                                    {{ end -}}
                                    {{ range .Categories -}}
                                    <a href="/{{.}}.html" class="badge badge-light" alt="{{.}}">{{.}}</a>
                                    {{ end -}}
                                    {{ range .Tags -}}
                                    <a href="/{{.}}.html" class="badge badge-light" alt="{{.}}">{{.}}</a>
                                    {{ end -}}
                                </div>
                                <div class="card-footer bg-dark text-white">
                                    <small>{{ .PublishDate }}</small>
                                    <small class="float-right">{{ .Duration }}</small>
                                </div>
                            </div>
                        </div>
{{ end -}}