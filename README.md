# go.demo.iam
Building an IAM service - tutorial project


This tutorial project shows how to build an IAM service including user login (SSO).
You can follow this [playlist](https://www.youtube.com/watch?v=Lg5LCdYU3Go&list=PLk2imHwYa7D51PlAZ-xMq7LtWI-rQ4vr4) 
and checkout the branch that is mentioned under each video to properly follow the tutorial.

Since this is a demo, I maybe did not respect all "guidelines" for having good software quality (e.g. most of the config
is hardcoded). If you want to build something on top of this, I strongly recommend changing that along with checking everything
for proper security.

You can find the library for validating the tokens issued by this IAM service [here](https://github.com/dtop/go.demo.jwt.lib)

## Installation

1)

```

$ cd <your go path>/src/github.com
$ mkdir dtop
$ cd dtop/
$ git clone git@github.com:dtop/go.demo.iam
$ cd go.demo.iam/
$ git checkout <branch> 
$ glide install

```

if you dont have glide, get it here: [Masterminds/glide](https://github.com/Masterminds/glide)

## Credits

I've used a couple of cool libraries to make this happen, just take a look into the glide.yaml to
see what is used.
