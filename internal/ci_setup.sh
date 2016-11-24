$GOBIN/natsc request -s $NATS_URI -t 5 -r 99 'group.set' '{"id":"1","name": "ci_admin"}'
$GOBIN/natsc request -s $NATS_URI -t 5 -r 99 'user.set' '{"group_id": 1, "username": "ci_admin", "password": "pwd", "admin":true}'
ernest-cli target $CURRENT_INSTANCE
ernest-cli login --user ci_admin --password pwd
ernest-cli user create usr pwd
ernest-cli group create test
ernest-cli group add-user usr test
ernest-cli login --user usr --password pwd
ernest-cli datacenter create vcloud fake --vcloud-url https://myvdc.me.com --fake --user usr --password pwd --org org --vse-url http://localhost --public-network NETWORK
ernest-cli datacenter create aws fakeaws --region fake --secret_access_key my_very_large_access_key_up_to_16 --access_key_id secret_key_up_to_16_chars --fake
