- [ ] Project 2 (SEC1): [[Project 2 - Description (SEC1)]] [[Project 2 - Tips (SEC1)]] 📅 2023-10-31 
## Problem
- We don't want anyone to know the data
- We don't want anyone to modify the data

## Adverseries
Dolev-yao
Recipients can be adversaries
	- semi-honest parties or passive adversaries
We don't have honest parties
	- They don't trust each other to see or process the data.
	- They only trust each other to follow the protocol they have agreed in collecting the data.
	- This means we have a semi-honest or passive adversary
	- We don't trust the hospital with the data and we don't trust the other patients with our own data

## Solution: Secure aggregation
Instead of looking at all patients individually we aggregate all the patients values into one single value
- Can be done by simply summing all the parties values into one
	- We can still learn what their total values are less than
		- This should be fine though
- We can do it with MPC protocol as it can sum the parties values
	- We can then send this to the hospital
	- Instead of implementing secure problem we can use the simple implementation in [[Lecture 6 - Slides (SEC1).pdf]]
	- We also need integrity as the dolev-yao can modify the message
	- We can't send the sum to the hospital in the clear, so we still have to us TLS
- Don't implement TLS yourself. Never implement cryptography yourself.
- We can use integers with modulo a prime for our field representation
	- We can whatever we want, but recommendation is to use integers with modulo a prime
In the report we should describe
- The problem
- The adverseries
- We have to describe our solution
- and why it solves our problem
Then implement it in code

It is ok to make sockets on the local host, we don't have to use different clients using vms ect.
It is ok to just accept certificates that a self signed, as it is not part of the learning goal


