import str from '../snpStore';
import {LocalStorageMock, LocationMock} from "../../test/browserMocks";
import * as util from "../../site-wide/util";
import {uri} from "../../site-wide/util";
import * as cl from '../../api/client';
import assert from 'assert';

cl.httpPost = jest.fn(() => {
    return new Promise(res => {
        res({});
    });
});

cl.httpPut = jest.fn(() => {
    return new Promise(res => {
        res({});
    });
});

const draft = {al:{"n": "draft", "p": "test-pouch", "u": "richard", "x": "txt"}};
const draftInStorage = {al:{"draft" : true, "n": "draft", "p": "test-pouch", "u": "richard", "x": "txt"}};

describe('snp store', () => {

    beforeAll(() => {
        global.localStorage = new LocalStorageMock();
    });

    test('Creating', (done) => {
        str.currPch = {p: {p: 'test-pouch'}, items: []};
        str.createDraft('richard').then(() => {
            let snp = str.currSnp;

            assert.deepEqual(snp.al, draftInStorage.al);
            assert.equal(snp.al.draft, true);
            assert.deepEqual(str.currPch.items[0], snp);

            // drafts
            let ds = str.getDrafts(snp.al.u, snp.al.p);
            assert.deepEqual(ds[0], snp);

            // get override
            let snp2 = str.getSnp(uri(snp.al));
            assert.deepEqual(snp2, snp);
            str.currSnp.cnt = 'some work in progress';
            str.currSnp.al.n = 'a new name';
            return;
        }).then(() => {
            return str.updSnp();
        }).then(() => {
            assert.equal(cl.httpPost.mock.calls[0][0], 'snippets');
            assert.deepEqual(cl.httpPost.mock.calls[0][1], str.currSnp);
            assert.equal(str.currSnp.al.draft, undefined);
            let ds = str.getDrafts('richard', 'test-pouch');
            assert.equal(str.currPch.items[0].al.draft, undefined);
            assert.equal(ds, null);

            let fsnp = str.getSnp(uri(str.currSnp.al));
            assert.equal(fsnp.cnt, str.currSnp.cnt);

            done();
        });
    });

    test('Should delete selected', (done) => {
        let pn = util.pathname;
        util.pathname = function () {
            return '/richard/test-pouch/draft.txt';
        };
        assert.deepEqual(str.getSelected(), [draft.al]);
        str.deleteSnippets().then(res => {
            assert.equal(cl.httpPut.mock.calls[0][0], 'snippets-delete');
            assert.deepEqual(cl.httpPut.mock.calls[0][1], { als: [draft.al]});
            done();
            util.pathname = pn;
        });
    });

    test('Should delete draft', (done) => {
        let pn = util.pathname;
        util.pathname = function () {
            return '/richard/test-pouch/draft.txt';
        };
        str.setDraft('richard', 'test-pouch', draft);
        str.deleteSnippets().then(res => {
            assert.equal(cl.httpPut.mock.calls[0][0], 'snippets-delete');
            assert.deepEqual(cl.httpPut.mock.calls[0][1], { als: [draft.al]});
            done();
            util.pathname = pn;
        });
    });
});