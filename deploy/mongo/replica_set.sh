#/bin/bash

mongosh --eval "rs.initiate({ _id: 'rs0', members: [{_id: 1, host: 'mongo:27017'}]});"


