# lab

[![godoc reference](https://cdn.rawgit.com/mvader/2faf5060e6cb109617ef5548836532aa/raw/2f5e2f2e934f6dde4ec4652ff0ae6d5c83cbfd6a/godoc.svg)](https://godoc.org/github.com/mvader/lab) [![Build Status](https://travis-ci.org/mvader/lab.svg?branch=master)](https://travis-ci.org/mvader/lab) [![codecov](https://codecov.io/gh/mvader/lab/branch/master/graph/badge.svg)](https://codecov.io/gh/mvader/lab)  [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)

![Gopher doing science](http://i.imgur.com/a0AMGoS.jpg)

`lab` is a library to handle feature switches and experiments with users.

### Getting started

#### The lab

First step is to create your lab to do science.
```go
myLab := lab.New()
```

#### Experiments

Then you can define your experiments.
```go
myLab.Experiment("new-footer-design").
	Aim(lab.AimPercent(60))

myLab.Experiment("show-website-upside-down").
	Aim(lab.AimRandom())
```

Experiments have aims, which are the audiences of the experiment, that is, the people that will be able to see the experiment.
These are the following kinds of aims:

* **Percent:** a percentage of the users will see the experiment.
* **Random:** users will randomly see the experiment.
* **Strategy:** users will see the experiment if they match the criteria of a defined strategy.
* **Everyone:** shows the experiment to everyone.
* **Nobody:** show the experiment to nobody.

You can use logical or and logical and to combine aims.

```go
myLab.Experiment("new-footer-design").
	Aim(lab.Or(
		lab.AimPercent(60),
		lab.And(
			lab.AimStrategy("is-admin", nil),
			lab.AimStrategy("is-not-john", nil),
		),
	))
```

Experiment `Aim` are chainables and act like logical ands.

```go
myLab.Experiment("john-birthday").
	Aim(lab.AimStrategy("is-admin", nil)).
	Aim(lab.AimStrategy("is-not-john", nil))
```

#### Strategies

You can define your own strategies. Strategies are used by aims to know if the experiments should be shown or not to visitors.

```go
myLab.DefineStrategy("visitor-is-the-beast", func(v lab.Visitor, p lab.Params) bool {
	return v.ID() == "666"
})
```

#### The session

In order to run the experiment you need to create a session. A session is a lab playground for a particular user.

```go
session := myLab.Session(myVisitor)
```

Visitors need to implement the `lab.Visitor` interface, which only requires to implement a method returning the string ID of the visitor.

#### Launch experiments

Finally, with the session, you can launch the experiments.

```go
wasRun := session.Launch("visitor-is-the-beast", nil)
```

If the callback is `nil`, `Launch` will only report if the experiment should be shown. If the callback is given, it will be run.

```go
session.Launch("visitor-is-the-beast", func() {
	wreakHavoc()
})
```

### Example

You can see a silly small web application demo to learn how to use the library in the example folder.

### Credits

Inspired by AirBnb's [trebuchet](https://github.com/airbnb/trebuchet).
