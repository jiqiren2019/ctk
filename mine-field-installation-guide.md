## system requirement
linux OS , Suggested Use centos 
Download the program from the official website first，Then buy the certificate and create the wallet.
When you're done, modify the configuration file and run


## Profile description
```
miner = Mine revenue wallet address

cert_id = Mine income wallet address mine certificate key 

# Listening port number
listen_port =  7555

[child]

# [child.xx]
# xx data format is digital
# The first trust miner USES 1,The second trustee miner USES 2,and the like


# First trust miner
[child.1]
miner = Miner income wallet address
cert_id = Miner income certificate key

# Second trust miner
[child.2]
miner = ....
cert_id = ..


# A third trust miner
[child.3]
miner = ....
cert_id = ..

# More trust miners
```


## Notes
Because mining program needs to transmit data before the node, only external network direct connection and port transmission are supported at present.
If it is running under a router or firewall, {mine program monitoring port number} should be used for data forwarding.
