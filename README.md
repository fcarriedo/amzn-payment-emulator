## Amazon Payment Flow Emulator

Quick & dirty Amazon Payment Flow emulator. Created to test some ideas on
payment integration using the co-branded flow. This may also apply to PayPal
since they provide similar integration.

Some basic principles of [webhooks](http://www.webhooks.org/) are used.

### Building it

Easy as a pie (it is obvious that I have not idea as to how hard is to cook a
pie)!

```
  $ git clone https://github.com/fcarriedo/amzn-payment-emulator.git && cd amzn-payment-emulator
  $ go get
  $ go build -o amzn
```

### Running the Server

```
Usage:  ./amzn [options]

Starts the Amazon/Paypal emulation server

  -p=9500: the port to listen for connections
```

### Using it

The entry point URL expects to have been redirected with several URL parameters
that describe the payment authorization as well as the callback URL.

Example:

  1. after running the app (`./amzn -p 8080`), point your browser to the
     following example URL:
`http://localhost:8080/auth?callbackURL=http%3A%2F%2Fes-razrzone.ngrok.com%2Fcallback&amount=5&desc=Cosmos%20A%20Personal%20Voyage%20DVD%20Set&id=aaa123`
  1. Authenticate to your imaginary account (any user/pwsd would do)
  1. Check your order and confirm the order (the given callback URL expects to
     be available or it will pop an error upon confirmation).

### License

Of course covered by [WTFPL License](http://www.wtfpl.net/).
