cd /usr/local/services/marketing_api_tool-1.0/bin/
for ((i=-30; i<=-3; i++))
do
echo $i >> ../log/marketing_api_tool_30_days_ago.log
./marketing_api_tool -datetime=$i -sleep=true -condition="STAT_GROUP_BY_TIME_DAY,STAT_GROUP_BY_AD_ID" >> ../log/marketing_api_tool_30_days_ago.log
./marketing_api_tool -datetime=$i -sleep=true -condition="STAT_GROUP_BY_TIME_DAY,STAT_GROUP_BY_CREATIVE_ID,STAT_GROUP_BY_INVENTORY" >> ../log/marketing_api_tool_30_days_ago.log
done
