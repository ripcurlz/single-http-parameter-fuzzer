# single http parameter fuzzer

**DISCLAIMER: ONLY USE THIS PROGRAM ON TARGETS YOU HAVE PERMISSION TO FUZZ!**      

Initially used as a "poor man's" http fuzzer for one http parameter instead of the very slow "Intruder" plugin in burpsuite community, used for the first time to fuzz on the tryhackme room "sustah" with an upcounting number in between a given range and to check if a given string is or is not in the response body (go check it out: https://tryhackme.com/room/sustah :) ). I know there is probably something out there I could have used instead, but I also did it to learn some golang along the way :)

```
Usage (important: the parameter "-headers" has currently to be the last one of the flags if you use it!):
  -host string
        the IP address or hostname you want to fuzz on
  -port int
        the port you want to use (default 80)
  -method string
        the http method you want to use, either 'get' or 'post'
  -parametertofuzz string
        the single http parameter you want to fuzz
  -stringnottobeinresponse string (default "")
        the string you do not want to be in the http response body for a match
  -stringtobeinresponse string (default "")
        the string you do want to be in the http response body for a match
  -wordlist string
        path of the wordlist you want to use for fuzzing the parameter, if you do not use startnumber and endnumber
  -startnumber int
        the start number to use for fuzzing the parameter (default 0)
  -endnumber int
        the end number to use for fuzzing the parameter (default 100)
  -headers string (optional)
        http headers you want to use, in the form of 'header1:value1' 'header2:value2' and so on
```
IMPORTANT:
You can EITHER use the flag "-wordlist" or INSTEAD both flags "-startnumber" and "-endnumber" for fuzzing.

example with startnumber 10900 and endnumber 99999:

```
go run main.go -method post -host 10.10.63.138 -port 8085 -stringnottobeinresponse "Oh no" -parametertofuzz number -startnumber 10900 -endnumber 99999 -headers "Content-Type:application/x-www-form-urlencoded" "X-Originating-IP:127.0.0.1" "X-Forwarded-For:127.0.0.1" "X-Remote-IP:127.0.0.1" "X-Remote-Addr:127.0.0.1" "X-Client-IP:127.0.0.1" "X-Host:127.0.0.1" "X-Forwarded-Host:127.0.0.1"


```


example with wordlist "./wordlist":

```
go run main.go -method post -host 10.10.63.138 -port 8085 -stringtobeinresponse "lucky" -parametertofuzz number -wordlist ./wordlist -headers "Content-Type:application/x-www-form-urlencoded" "X-Originating-IP:127.0.0.1" "X-Forwarded-For:127.0.0.1" "X-Remote-IP:127.0.0.1" "X-Remote-Addr:127.0.0.1" "X-Client-IP:127.0.0.1" "X-Host:127.0.0.1" "X-Forwarded-Host:127.0.0.1"


```