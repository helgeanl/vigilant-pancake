// Testfile
package driver

var aName = 2// private

var BigBro =3// public (exported)

//var 123abc // illegal
/*
 func (p *Person) SetEmail(email string) {  // public because SetEmail() function starts with upper case
  	p.email = email
 }
*/
 func Test() int { // private because email() function starts with lower case
  	return 10
 }

 func test() int { // private because email() function starts with lower case
   return 5
 }
