/* 加密 */
export async function encrypt(key, plaintext) {
    try {
        // 将密钥转换为 ArrayBuffer
        const keyBuffer = new TextEncoder().encode(key);

        // 导入密钥
        const importedKey = await window.crypto.subtle.importKey(
            "raw", keyBuffer, { name: "AES-GCM" }, false, ["encrypt"]
        );

        // 创建一个随机的 nonce
        const nonceSize = 12; // AES-GCM 默认的 nonce 长度是 12 字节
        const nonce = crypto.getRandomValues(new Uint8Array(nonceSize));

        // 将 plaintext 转换为 Uint8Array
        const plaintextBuffer = new TextEncoder().encode(plaintext);

        // 使用 AES-GCM 加密
        const ciphertext = await window.crypto.subtle.encrypt(
            {
                name: "AES-GCM",
                iv: nonce, // 使用生成的 nonce
            },
            importedKey,
            plaintextBuffer
        );

        // 将 ciphertext 转换为 base64 URL 编码格式
        const ciphertextArray = new Uint8Array(ciphertext);
        let base64Url = btoa(String.fromCharCode(...ciphertextArray));

        // 将 = 替换为空字符串，以符合 base64url 编码的规范
        base64Url = base64Url.replace(/=+$/, "");

        // 将 nonce 和 ciphertext 结合起来（nonce 作为前缀）
        const nonceBase64Url = btoa(String.fromCharCode(...nonce)).replace(/=+$/, "");
        let result = nonceBase64Url + base64Url;
        result = result.replace(/\+/g, '-').replace(/\//g, '_');
        return result;
    } catch (error) {
        console.error("Encryption failed:", error);
        throw error;
    }
}

/* 解密 */
export async function decrypt(key, encryptedText) {
    // 将 Base64 URL 编码转换为标准 Base64 编码
    encryptedText = encryptedText.replace(/-/g, '+').replace(/_/g, '/');

    // 添加 padding
    const paddingLength = encryptedText.length % 4;
    if (paddingLength !== 0) {
        encryptedText += "=".repeat(4 - paddingLength);
    }

    try {
        // 将密钥转换为 ArrayBuffer
        const keyBuffer = new TextEncoder().encode(key); // 使用 TextEncoder 编码为字节数组

        // 导入密钥
        const importedKey = await window.crypto.subtle.importKey(
            "raw", keyBuffer, { name: "AES-GCM" }, false, ["decrypt"]
        );

        // 解码 base64 字符串
        const decodedData = Uint8Array.from(atob(encryptedText), c => c.charCodeAt(0));

        const nonceSize = 12; // AES-GCM 默认的 nonce 长度是 12 字节
        if (decodedData.length < nonceSize) {
            throw new Error("Ciphertext is too short");
        }

        // 提取 nonce 和 ciphertext
        const nonce = decodedData.slice(0, nonceSize);
        const ciphertext = decodedData.slice(nonceSize);

        // 解密过程
        const decryptedData = await window.crypto.subtle.decrypt(
            {
                name: "AES-GCM",
                iv: nonce,
            },
            importedKey,
            ciphertext
        );

        // 解密结果转换为字符串
        return new TextDecoder().decode(decryptedData);
    } catch (error) {
        console.error("Decryption failed:", error);
        throw error;
    }
}