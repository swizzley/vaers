# Overview
An attempt to catalog and dive deeper into the VAERS data. I downloaded it 1 month at a time because 
it was giving me too many rows to return errors, max is 10,000. All I was trying to do is get DEATH from the 
Covid-19 Vaccine from 2021 and depending on the query it was between 20-50,000 results. So I split up the query 
by date, combined them all, and ended up with ~5000 deaths which is less than if you just download it all as CSV. 

The data seems to change constantly, and I mean it decrements, which should be impossible. I had 1 query that was 
returning ~21,000 results, and then 5min later it was return 18,000 results. This combined data should have been 
all results, but it is clearly a tiny subset. 

In my second attempt at cataloging the data directly from the data download feed, the real data problems begin to emerge.
Events where symptoms were DEATH, but not marked as a death, and or have no death date was one of the most obvious.
The second was there are obviously typos, in that some people are marked to have died in 2001, but they were reported 
in 2021. 

# Conclusion
  * The data is still very bad
  * Not all deaths are labeled deaths
  * Not all deaths have death dates
  * There is no consistent date value at all, meaning there is no single date field that is valid among all events
  * Miscarriage is not marked as death or birth defect 
  * There are 2 age fields and neither is always present 
  * Lots and lots of null values everywhere 
  * Thousands of unknown sex 
  * Seemingly invalid state Codes `XB` and `GU` 
  * **THOUSANDS OF DEATHS LABELED AS A SYMPTOM INSTEAD OF LABELING REACTION AS DEATH**
  
![visualization](dashboard.png "Title")


# Usage

  1. open a browser and download the capcha protected file from [vaers](https://vaers.hhs.gov/data/datasets.html?).
  1. Unzip the file to `./data/`
  1. ```make run```
  1. ```cd vaers-data-loader && go build -o vaers-data-loader```
  1. `./vaers-data-loader`
  1. open [localhost:5601](http://localhost:5601)
  1. navigate in the menu to elastic -> Manage -> Management menu -> saved objects -> import -> import -> select the dashboard file in the vaers directory `vaers/dashboard-objects.ndjson` use default settings and click Import
  1. navigate to the dashboard or begin searching the data in the discover section
