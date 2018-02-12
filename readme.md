## Jumpcloud Challenge: Mini Golang Api

#### Quick Backstory
I've never actually written golang before this.  I spent yesterday (Saturday) reading documentation (gobyexample.com was indispensable) and then today (Sunday) building the sample application.


### Server
Uses port `:8000` and can be run using `go run main.go`

I used
```
for i in `seq 1 2000`; do curl --data "password=angryMonkey" http://localhost:8000/hash; done;
```
to test.


### Notes on implementation

I haven't figured out file structure best practices for golang so to keep this simple everything is in `main.go`.  At each step I copied my main into solution files labeled accordingly.  Although I was able to solve steps 3-5 in one step because of how related they are.

I leaned into the more verbose side of logging to document the important steps at each method call.  Hopefully this helps the results verification process.  If this was actually production, I'd replace most of those console logs with proper test coverage to ensure that at each functional step my business logic was still sound.

The requested precision of milliseconds is not precise enough to show meaningful metrics because this job completes in .4-.7 milliseconds.  When the artificial delay of 5 seconds is introduced, the average always comes out to be 5000 (without the delay it rounds milliseconds down to 0)  A different precision would need to be implemented to improve this number.

My use of a Map construct is a simple replacement for an external data store that would persist storage across multiple instances of the service in production.
