cd /usr/local/services/marketing_api_tool-1.0/bin/
./marketing_api_tool -datetime=0 -sleep=false -condition="STAT_GROUP_BY_TIME_HOUR,STAT_GROUP_BY_AD_ID" >> ../log/marketing_api_tool.log
./marketing_api_tool -datetime=0 -sleep=false -condition="STAT_GROUP_BY_TIME_HOUR,STAT_GROUP_BY_CREATIVE_ID,STAT_GROUP_BY_INVENTORY" >> ../log/marketing_api_tool.log
