cd /usr/local/services/marketing_api_tool-1.0/bin/
./marketing_api_tool -datetime=-2 -sleep=true -condition="STAT_GROUP_BY_TIME_DAY,STAT_GROUP_BY_AD_ID" >> ../log/marketing_api_tool_two_days_ago.log
./marketing_api_tool -datetime=-2 -sleep=true -condition="STAT_GROUP_BY_TIME_DAY,STAT_GROUP_BY_CREATIVE_ID,STAT_GROUP_BY_INVENTORY" >> ../log/marketing_api_tool_two_days_ago.log
