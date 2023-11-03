const express = require('express');
const bodyParser = require('body-parser');
const applySteps = require('./prosemirror-node');
const app = express();
const port = 8080;

app.use(bodyParser.json());
app.post('/', (req, res) => {
  const payload = req.body;
  const doc = applySteps(payload.doc, payload.steps)
  res.status(200).send(JSON.stringify({ doc }));
});

app.listen(port, () => {
  console.log(`node server listening...`);
});
