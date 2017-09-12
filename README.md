# File Storage Microservice

A small HTTP application running in docker that accepts files of various styles, does it's best to verify they are what they claim to be (including clam av) before storing the specified storage repository. Files or S3 for now.

The system is intended to abstract away a lot of common file storage functions carried out by web apps and provide a stateless, scalable file storage mechanism that can be used across multiple projects.
