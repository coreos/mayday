# mayday [![Build Status](https://travis-ci.org/coreos/mayday.png?branch=master)](https://travis-ci.org/coreos/mayday)

...man overboard!

```


     ___ ___   ____  __ __  ___     ____  __ __ 
    |   |   | /    ||  |  ||   \   /    ||  |  |
    | _   _ ||  o  ||  |  ||    \ |  o  ||  |  |
    |  \_/  ||     ||  ~  ||  D  ||     ||  ~  |
    |   |   ||  _  ||___, ||     ||  _  ||___, |
    |   |   ||  |  ||     ||     ||  |  ||     |
    |___|___||__|__||____/ |_____||__|__||____/ 
                                                



~~~~~~~~~~~~~~~~~~~~~~~~~lol~~~~~~~~~~~~~~~~~~~~~~~~

                                   ><>
               <><  
```

## overview
Mayday is a tool to simplify gathering support information.  It is built in the
spirit of sysreport, son of sysreport (sosreport), and similar support tools.  
Mayday gathers information about the configuration, hardware, and running state
of a system.


## goals
The goals of mayday are:
  * simplify gathering information about a running system into a single command
  * collect information into one single file to be transferred to support staff
  * when possible the file should be small enough to be sent via email (<10MB)
  * *not* collect sensitive information like crypto keys, password hashes, etc
  * extensible through plugin system
