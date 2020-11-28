    <body>
        <header>
            <div class="collapse bg-white" id="navbarHeader">
                <div class="container">
                <div class="row">
                    <div class="col-sm-8 col-md-7 py-4">
                    <h4 class="text-dark">About</h4>
                    <p class="text-muted">This is a personal project, developed in Go, using BeeGo. I hope you enjoy it!</p>
                    <form action="https://www.paypal.com/cgi-bin/webscr" method="post" target="_top">
                        <p class="text-muted">This is an ad-free project. If you like it, please consider a donation to help me cover the costs of keeping the project online. Thank you!</p>
                        <input type="hidden" name="cmd" value="_s-xclick" />
                        <input type="hidden" name="hosted_button_id" value="PE2N8QAY7N28G" />
                        <input type="image" src="https://www.paypalobjects.com/es_ES/ES/i/btn/btn_donate_LG.gif" border="0" name="submit" title="PayPal - The safer, easier way to pay online!" alt="Botón Donar con PayPal" />
                        <img alt="" border="0" src="https://www.paypal.com/es_ES/i/scr/pixel.gif" width="1" height="1" />
                    </form>
                    </div>
                    <div class="col-sm-4 offset-md-1 py-4">
                    <h4 class="text-dark">Contact</h4>
                    <ul class="list-unstyled text-muted">
                        <li><a href="https://twitter.com/dsmatilla" class="text-dark">Follow author on Twitter</a></li>
                        <li><a href="https://github.com/dsmatilla" class="text-dark">My github</a></li>
                        <li><a href="mailto:daniel@esdis.es" class="text-dark">Email me</a></li>
                    </ul>
                    </div>
                </div>
                </div>
            </div>
            <div class="navbar navbar-light box-shadow">
                <div class="container d-flex justify-content-between">
                <a href="/" class="navbar-brand d-flex align-items-center">
                    <img src="/img/just-tit.png" alt="Just-tit logo" width="150px" height="50px" />
                </a>
                <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarHeader" aria-controls="navbarHeader" aria-expanded="false" aria-label="Toggle navigation">
                    <span class="navbar-toggler-icon"></span>
                </button>
                </div>
            </div>
        </header>

        <main role="main">
            <section class="jumbotron text-center bg-white">
                <div class="container">
                    <div class="col-md-10 col-lg-8 col-xl-7 mx-auto">
                    <form action="/" method="get">
                        <div class="form-row">
                        <div class="col-12 col-md-9 mb-2 mb-md-0">
                            <input class="form-control form-control-lg" title="Search" id="search" type="text" name="s" placeholder="Enter your search...">
                        </div>
                        <div class="col-12 col-md-3">
                            <button type="submit" value="Search" class="btn btn-block btn-lg btn-primary">Search</button>
                        </div>
                        </div>
                    </form>
                    </div>
                </div>
            </section>

            <div class="album py-5 bg-white min-vh-100">
                <div class="container">
                    <div class="row">
                        {{ .LayoutContent }}
                    </div>
                </div>
            </div>
        </main>

        <footer class="text-muted bg-dark">
            <div class="container text-center">
                <img loading="lazy" src="/img/pornhub.png" alt="Pornhub logo" width="150px" height="50px" />
                <img loading="lazy" src="/img/redtube.png" alt="Redtube logo" width="150px" height="50px" />
                <img loading="lazy" src="/img/youporn.png" alt="Youporn logo" width="150px" height="50px" />
                <img loading="lazy" src="/img/tube8.png" alt="Tube8 logo" width="150px" height="50px" />
                <br />
                <img loading="lazy" src="/img/spankwire.png" alt="Spankwire logo" width="150px" height="50px" />
                <img loading="lazy" src="/img/xtube.png" alt="Xtube logo" width="150px" height="50px" />
                <img loading="lazy" src="/img/extremetube.png" alt="Extremetube logo" width="150px" height="50px" />
                <img loading="lazy" src="/img/keezmovies.png" alt="Keezmovies logo" width="150px" height="50px" />
                <div class="footer-copyright text-center py-3 text-white">© 2010 - 2020 Copyright:
                    <a class="text-white" href="https://esdis.cloud"> Esdis Cloud</a>
                </div>
            </div>
        </footer>
        
        <!-- Global site tag (gtag.js) - Google Analytics -->
        <script async src="https://www.googletagmanager.com/gtag/js?id=UA-106943798-2"></script>
        <script>
            window.dataLayer = window.dataLayer || [];

            function gtag() {
                dataLayer.push(arguments);
            }

            gtag('js', new Date());
            gtag('config', 'UA-106943798-2');
        </script>
        <script type="text/javascript">var sc_project=11999640;var sc_invisible=1;var sc_security="4e8e7b5a";</script>
        <script type="text/javascript" src="https://www.statcounter.com/counter/counter.js" async></script>
        <script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
        <script src="https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ho+j7jyWK8fNQe+A12Hb8AhRq26LrZ/JpcUGGOn+Y7RsweNrtN/tE3MoK7ZeZDyx" crossorigin="anonymous"></script>
        <script>
            if ('serviceWorker' in navigator) {
                navigator.serviceWorker.register('/service-worker.js');
            }
        </script>
        <script type="text/javascript" src="https://www.hubtraffic.com/js/external/helpers.js"></script>
        {{ if .Result -}}
            {{ if eq (index .Result 0).Type  "single" -}}
            <script type="application/ld+json">
                {
                    "@context": "https://schema.org",
                    "@type": "VideoObject",
                    "name": "{{(index .Result 0).Title}}",
                    "description": "{{(index .Result 0).Title}}",
                    "thumbnailUrl": [
                        "{{(index .Result 0).Domain}}{{ToImageProxy (index .Result 0).Thumb}}"
                    ],
                    "uploadDate": "{{ (index .Result 0).PublishDate }}",
                    "embedUrl": "{{(index .Result 0).URL}}?tp=true"
                }
            </script>
            {{ end -}}
        {{ end -}}
        <script type="application/ld+json">
            {
                "@context": "https://schema.org",
                "@type": "WebSite",
                "url": "https://just-tit.com",
                "potentialAction": {
                    "@type": "SearchAction",
                    "target": "https://just-tit.com/{search_term_string}.html",
                    "query-input": "required name=search_term_string"
                }
            }
        </script>
    </body>