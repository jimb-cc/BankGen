# BankGen
A GO script for generating random bank customer records. 
## Usage
`go run bankgen.go --help` for options
## Output
depending on what I've fixed or changed the output will look something like this:
```
{
	"_id" : ObjectId("5c0839dbaba8551adb818d82"),
	"name" : "Jacinthe Roob",
	"branch" : "Lights shire",
	"branch_id" : "EN",
	"manager" : "Dock Schimmel",
	"country" : "EN",
	"rankLevel" : NumberLong(5),
	"accounts" : [
		{
			"accountType" : "Current",
			"accountSubType" : "digitalCurrentAccount",
			"overdraftLimit" : NumberLong(1000),
			"balance" : 3413.48
		},
		{
			"accountType" : "Savings",
			"accountSubType" : "SuperSaver",
			"interestRate" : 1.514,
			"balance" : 22140.5
		},
		{
			"accountType" : "ISA",
			"accountSubType" : "SuperTaxFreeISA",
			"interestRate" : 3.208,
			"balance" : 33714.25
		},
		{
			"accountType" : "Mortgage",
			"accountSubType" : "BuildingDeluxe",
			"interestRate" : 4.691,
			"balance" : -179200.36
		}
	]
}
```
