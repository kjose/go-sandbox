#!/bin/bash
GOOS=linux GOARCH=amd64 go build -o go-service
scp -i ~/.ssh/kp-20201204.pem go-service ubuntu@ec2-35-180-47-155.eu-west-3.compute.amazonaws.com:
scp -i ~/.ssh/kp-20201204.pem -r templates ubuntu@ec2-35-180-47-155.eu-west-3.compute.amazonaws.com:
scp -i ~/.ssh/kp-20201204.pem -r assets ubuntu@ec2-35-180-47-155.eu-west-3.compute.amazonaws.com:
ssh -i ~/.ssh/kp-20201204.pem ubuntu@ec2-35-180-47-155.eu-west-3.compute.amazonaws.com "sudo chmod 700 go-service;
    sudo systemctl stop go-sandbox.service;
    sudo touch /etc/systemd/system/go-sandbox.service;
    echo '[Unit]
Description=Go server

[Service]
ExecStart=/home/ubuntu/go-service
WorkingDirectory=/home/ubuntu
User=root
Group=root
Restart=always

[Install]
WantedBy=multi-user.target' | sudo tee /etc/systemd/system/go-sandbox.service;
    sudo systemctl enable go-sandbox.service;
    sudo systemctl start go-sandbox.service;"
