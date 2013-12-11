litr
====
litr.io

# Stack
* Go
	* [Martini](http://martini.codegangsta.io/)
        * [gorp](http://shadynasty.biz/blog/2012/09/05/auth-and-sessions/)
* Postgres
	* Heroku Postgres
	* pq (go driver)
	* Induction (database client)
* Heroku
* S3 + CloudFront
	* [goamz](https://github.com/crowdmob/goamz)
* wercker for deployment? (later)

# Reference
## Heroku
* `heroku logs`

## PostgreSQL
* `heroku pg:info`
	* db info
* `heroku pg:psql`
	* db session
## Go
### Tutorials
* [Wiki](http://golang.org/doc/articles/wiki/)
* [Guestbook](http://shadynasty.biz/blog/2012/07/30/quick-and-clean-in-go/)
* [Getting Started With Go On Heroku](http://mmcgrana.github.io/2012/09/getting-started-with-go-on-heroku.html)
* [Auth and Sessions](http://shadynasty.biz/blog/2012/09/05/auth-and-sessions/)
* [gorp tutorial](http://nathanleclaire.com/blog/2013/11/04/want-to-work-with-databases-in-golang-lets-try-some-gorp/)

## XSS
* HTML sanitization

## CSRF

## SSL

## DNS
* limitations of DNS A-records
	* https://devcenter.heroku.com/articles/apex-domains
* Zerigo
	* https://devcenter.heroku.com/articles/zerigo_dns
* need ALIAS record
