const { Node, Schema } = require("prosemirror-model")
const { Step } = require("prosemirror-transform")

const schema = new Schema({
  nodes: {
    text: {},
    doc: {content: "text*"}
  }
})

module.exports = function applySteps(docJSON, stepsJSONS) {
	let doc = Node.fromJSON(schema, JSON.parse(docJSON));
	stepsJSONS = JSON.parse(stepsJSONS)
	stepsJSONS.forEach((s) => {
	  const step = Step.fromJSON(schema, s);
	  doc = step.apply(doc).doc;
	});
	return JSON.stringify(doc);
}
