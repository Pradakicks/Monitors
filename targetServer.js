const express = require('express');
const mongoose = require('mongoose');
const bodyParser = require('body-parser');
const cors = require('cors')
const items = require("./routes/api/items");
require ('newrelic');


const app = express();

// Bodyparser Middleware;
app.use(cors())
app.use(bodyParser.json());

// DB Config 
const db = require('./config/keys').mongoURI;

// Connect to Mongo
mongoose
    .connect(db, { useNewUrlParser: true, useUnifiedTopology: true })
    .then(() => console.log('MongoDB Connected...'))
    .catch(err => console.log(err));

// Use Routes
app.use('/api/items', items)
 

    const port = process.env.PORT || 5003;

    app.listen(port, () => console.log(`Server started on port ${port}`));