const { Node, Schema } = require("prosemirror-model")
const { Step } = require("prosemirror-transform")

const schema = new Schema({
  nodes: {
    text: {},
    doc: {content: "text*"}
  }
})

module.exports = function applySteps(inDoc, steps) {
	let doc = Node.fromJSON(schema, inDoc);
	steps.forEach((s) => {
	  const step = Step.fromJSON(schema, s);
	  doc = step.apply(doc).doc;
	});
	return doc;
}
