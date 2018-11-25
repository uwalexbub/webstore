# About

The focus of this project is:
* Learn Go language and how to build web applications with it
* Learn how to instrument Go code with metrics, collect and visualize them
* Explore Linux tools for monitoring and collecting system and application performance metrics and how to interpret them in the context of the application.
* Learn how to profile Go code and identify bottlenecks

This project implements a simple HTTP storage service and a client that runs performance tests against the service.

The service exposes HTTP APIs to:
* upload text with specified name
* download previously uploaded text by specified name
* 'clear' the service's internal storage: internally service saves texts as files

The service encrypts the text upon upload and decrypts upon download using simple symmetric encryption with a secret key.

The client runs a specified number of continuous tests in parallel, where each test:
* calls service to upload random text data of random sizes (within specified bounds) with randomly generated name
* calls service to download the uploaded text by specified name
* asserts that downloaded downloads the and calls upload and download HTTP APIs to test the service
* sleeps and repeats