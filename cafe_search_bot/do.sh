#!/bin/bash

gcloud builds submit --region=us-central1 --config=cloudbuild.yaml .
