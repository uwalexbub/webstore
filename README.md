# About

This project implements a simple HTTP storage service and a client that runs performance and stress tests against the service. The focus of this project is to learn how to instrument code, collect application and OS metrics, test a web application for performance and identify bottlenecks (CPU vs memory vs IO). 

The service allows to upload new files (or any binary blob) and download previously uploaded files.

The client generates random binary files of different sizes (1 MB, 10 MB and 100 MB) and tests the service.