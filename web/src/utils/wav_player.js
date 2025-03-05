// 播放 PCM16 语音流
export const playPCM16 = (pcm16Array, sampleRate = 44100) => {
    try {
        // 创建 AudioContext
        const audioContext = new (window.AudioContext || window.webkitAudioContext)();

        // 将 Int16Array 转换为 Float32Array (Web Audio API 使用 Float32)
        let float32Array = new Float32Array(pcm16Array.length);
        for (let i = 0; i < pcm16Array.length; i++) {
            float32Array[i] = pcm16Array[i] / 32768; // Int16 转换为 Float32
        }

        // 创建 AudioBuffer
        const audioBuffer = audioContext.createBuffer(1, float32Array.length, sampleRate); // 单声道
        audioBuffer.getChannelData(0).set(float32Array); // 设置音频数据

        // 创建 AudioBufferSourceNode 并播放音频
        const source = audioContext.createBufferSource();
        source.buffer = audioBuffer;
        source.connect(audioContext.destination); // 连接到扬声器
        source.start(); // 播放
        return source
    } catch (e) {
        console.warn(e)
        return null
    }
}
