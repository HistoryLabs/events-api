import express from 'express';

const app = express();

app.use('/', (req, res) => {
    res.send('This is the events-api! Welecome!');
});

app.listen(3000, () => console.log('events-api listening on port 3000'));