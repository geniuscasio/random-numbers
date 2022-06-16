### How to run
```
docker build . -t random-numbers
docker run -p 8080:8080 random-numbers
sh test.sh
```

### Some thoughts
* All limits described in [API doc](https://www.random.org/clients/http/) I've tried to cover.
* REST _always-200-rule_ - there is a big discussions about using HTTP codes in responses or not, but
  I decided to use 200 for everything, and showing errors as a JSON messages here.
* Sessions - probably would use something like 
  [github.com/gin-contrib/sessions](https://github.com/gin-contrib/sessions) for real world app
* Logging - I've tried to use [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) for Gin, 
  but it didn't go well, so I've skipped that to save time).
* Security - I think MD5 hashing that I've used is not best way to store passwords, adding salt would
  make it much better, also skipped for time reason.
* Tests - I've tried to use interfaces where it possibly, so we can create mocks from adapters using
  `mockery`, and then writing tests for `logic` package would be easy work.
* Performance - I bet sessions can work faster, but since `random.org` is quite slow, 
  don't think it even matters.
* I never used Gin before, everything is written in parallel with reading docs about it.
* * Hence above, I could miss some Gin best practices.

## P.S.
In the end I want to give some feedback form my end:)
I really missed few truly Go-ish things in the test:
* No room for concurrency ~~hell~~ fun;
* No room for using Go channels;