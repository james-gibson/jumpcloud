##Jumpcloud Challenge: Mini Golang Api

####Quick Backstory
I've never actually written golang before this.  I spent yesterday (Saturday) reading documentation (gobyexample.com was indispensable) and then today (Sunday) building the sample application.

As such I'd like to add a few notes, I haven't figured out file structure best practices for golang so to keep this simple everything is in `main.go`.  To document the solution at each step I copied my main into solution files labeled accordingly.  Although I was able to solve steps 3-5 in one step because of how related they are.

The second being I leaned into the more verbose side of logging to document the important steps at each method call.  Hopefully this helps the results verification process.  If this was actually Production I'd replace most of those console logs with proper test coverage to ensure that at each functional step my business logic was still sound.

My use of a Map construct is a simple replacement for an external data store that would persist storage across multiple instances of the service in production.

###Server
Uses port `:8000` and can be run using `go run main.go`

I used
```
for i in `seq 1 2000`; do curl --data "password=angryMonkey" http://localhost:8000/hash; done;
```
to test.


 
