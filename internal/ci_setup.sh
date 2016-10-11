$GOBIN/natsc request -s $NATS_URI -t 5 -r 99 'group.set' '{"id":"1","name": "ci_admin"}'
$GOBIN/natsc request -s $NATS_URI -t 5 -r 99 'user.set' '{"group_id": 1, "username": "ci_admin", "password": "pwd", "admin":true}'
