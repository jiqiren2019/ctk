## system requirement
linux OS , Suggested Use centos 
Download the program from the official website first，Then buy the certificate and create the wallet.
When you're done, modify the configuration file and run


## Profile description
```
# Enter mining account here(!!!!!!!!! Note: there is a space after the !!!!!!!!!!!)
# for example : 
# account: 0xb1f254c9b48b0681d9549b983c9f404692cd2a3b
account: 

# Enter the certificate code here(!!!!!!!!!!! Note: there is a space after the !!!!!!!!!!!!)
# for example : 
# ticket: 01234567
ticket: 

# Configure agent mining information
child:
  # Set generation mining account, multiple with, space
  # for example : 
  # -------------------------------------
  # account: [
  # 0xb1f254c9b48b0681d9549b983c9f404692cd2a3b,
  # 0xb1f254c9b48b0681d9549b983c9f404692cd2a3b,
  # ]
  # --------------------------------------
  account: [
  ]
  
  # Set generation mining certificate code, multiple with, separation
  # for example : 
  # -------------------------------------
  # ticket: [
  # 123456,
  # 7890ab,
  # ]
  # -------------------------------------- 
  ticket: [
  ]
```


## Notes
Because mining program needs to transmit data before the node, only external network direct connection and port transmission are supported at present.
If it is running under a router or firewall, {mine program monitoring port number} should be used for data forwarding.
