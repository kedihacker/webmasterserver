#!/bin/bash

sudo -s -H
apt-get clean
rm /var/lib/apt/lists/*
rm /var/lib/apt/lists/partial/*
apt-get clean
apt-get update