2 random shares fra vores field som er {0..10.000} (NOTE: Vi kan ikke bruge 10.000 irl, vi skal bruge en meget stor prime.)

#Step 1 -> Alle udregner deres chunks

##Alice
secret = vælger et tal
x1 = random tal mellem 0 .. 10.000
x2 = random tal mellem 0 .. 10.000
x3 = secret - ((x1 + x2) mod 10.000)

NOTE: (x1 + x2 + x3) mod 10.000 = secret (This would be how to find alice secret from the chunks)

##Bob
secret = vælger et tal
y1 = random tal mellem 0 .. 10.000
y2 = random tal mellem 0 .. 10.000
y3 = (secret - (y1 + y2)) mod 10.000

##Charlie
secret = vælger et tal
z1 = random tal mellem 0 .. 10.000
z2 = random tal mellem 0 .. 10.000
z3 = (secret - (z1 + z2)) mod 10.000

#Step 2 -> Alle deler deres secrets med hinanden

Alice: x1 x2 x3 -> alice bob charlie
Bob: y1 y2 y3 -> alice bob charlie
Charlie: z1 z2 z3 -> alice bob charlie


#Step 3 -> Alle har modtaget shares fra hinanden og udregner nu aggregated value

Alice: aliceAgg = (x2 + y3 + z1) mod 10.000
charlie: charlieAgg = (x1 + y2 + z3 ) mod 10.000	
bob: bobAgg = x3 + y1 + z2) mod 10.000


#Step 4 -> Alle har udregnet deres aggregated value og sender nu denne til hospitalet


#Step 5 -> Hospitalet modtager values fra bob, charlie og alice
Modtager: charlieAgg aliceAgg bobAgg

#step 6 -> Hospitalet udregner
Udregning
(charliAgg + aliceAgg + BobAgg) mod 10.000 eller charliAgg + aliceAgg + BobAgg
