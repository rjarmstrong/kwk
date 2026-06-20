import createElement from 'inferno-create-element';

function genVNodesFromRaw(raw) {
    let tags = {};
    let res = '';
    let readTag = false;
    let currTag = '';
    let closingTag = false;
    let tagNum = 0;
    for (let i = 0; i<raw.length; i++) {
        if (raw[i] === '<' && raw[i+1] === '/') {
            closingTag = true;
            readTag = false
        }
        if (raw[i] === '<') {
            readTag = true
        }
        if (raw[i] === '>') {
            if (closingTag) {
                // pop from states and create vnode
            } else {
                tags[currTag + (++tagNum)] = true
            }
            readTag = false
        }
        if (readTag) {
            currTag += raw[i]
        }
    }
    return res;
}