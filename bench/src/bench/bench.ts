import Axios from "axios";
import {JSEncrypt} from "jsencrypt";

const HOST_URL = "http://localhost:8080";

const MAX_ATTEMPTS = 500;

const MAX_KEYS = 100;

export function runBench() {

    let count = 0;

    const do_work:() => Promise<{}> = () => {
        return singleRun().then(() => {
            count++;

            if (count <= MAX_KEYS) {
                return do_work();
            }

            return;
        });
    };

    do_work().then(() => {
        console.log("GOOD");
        document.body.style.backgroundColor = "green";
    }, () => {
        console.log("ERROR");
        document.body.style.backgroundColor = "red";
    });
}

function singleRun() {
    let PUB_KEY = "";
    let UUID = "";
    let count = 0;

    return Axios.get(`${HOST_URL}/test/start`).then((data) => {
        PUB_KEY = data.data.public_key;
        UUID = data.data.test_id;

        const do_work:() => Promise<{}> = () => {
            return encryptEndVerify(PUB_KEY, UUID).then(() => {
                count++;

                if (count <= MAX_ATTEMPTS) {
                    return do_work();
                }

                return endTest(UUID);
            });
        };

        return do_work();

    });
}

function encryptEndVerify(pub_key:string, uuid:string):Promise<{}> {
    return Axios.get(HOST_URL + "/test/data/" + uuid).then((data) => {
        const encryptor = new JSEncrypt();

        encryptor.setPublicKey(pub_key);

        const encrypted = encryptor.encrypt(data.data.string);
        if (!encrypted) {
            return Promise.reject("fail to encrypt message");
        }

        return [encrypted, data.data.string];
    }).then(([encrypted, original]:[string, string]) => {
        return validate(encrypted, original, uuid);
    });
}

function validate(encrypted:string, original:string, uuid:string):Promise<{}> {
    return Axios.post(HOST_URL + "/test/verify/" + uuid, {
        encrypted,
        original

    }).then(() => {
        // asdsads;
        return {};
    });
}

function endTest(uuid:string):Promise<{}> {
    return Axios.get(HOST_URL + "/test/end/" + uuid).then(() => {
        // asdsads;
        return {};
    });
}
