# Title: DB to be used

# Status: Accepted

# Context
Our application needs a reliable and efficient database solution to store information about castles. MongoDB has been suggested as a potential option due to its flexibility, scalability, and suitability for storing JSON-like data structures.

# Decision
We have decided to use MongoDB as the database solution for storing castle data. This decision is based on bellow factors:

- Schema Flexibility: MongoDB's document-based model allows us to store data in a JSON-like format, making it easy to represent the hierarchical and semi-structured nature of castle information.

- Great free tier on Atlas: MongoDB Atlas offers a generous free tier with ample storage, bandwidth, and other resources, allowing us to start with minimal operational costs and scale as needed.

- Scalability: MongoDB is designed to scale horizontally, allowing us to handle large volumes of castle data and accommodate future growth without significant changes to our infrastructure.

- Querying Capabilities: MongoDB provides powerful querying capabilities, including support for complex queries, indexing, and aggregation pipelines, which are essential for retrieving and analyzing castle data efficiently.

- Geospatial Queries: MongoDB offers native support for geospatial queries, which will be useful for location-based searches and analysis of castle data - plans for long run ;).

- Community Support and Ecosystem: MongoDB has a large and active community, extensive documentation, and a rich ecosystem of tools and libraries that will facilitate development and maintenance tasks;

# Consequences
By choosing MongoDB as our database solution, we expect to benefit from its flexibility, scalability, and querying capabilities. 

Additionally, we need to monitor the performance of our MongoDB deployment and be prepared to scale our infrastructure as our data grows.

Overall, we believe that MongoDB is well-suited for storing castle data and will enable us to build a robust and scalable application.
