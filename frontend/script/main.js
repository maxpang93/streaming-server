import Stack from "./stack.js";

const BASE_URL = "http://localhost:8090"

const pageStack = new Stack("pageStack");

function setPreviousPage(previousPage) {
    pageStack.push(previousPage)
}

function getPreviousPage() {
    return pageStack.pop()
}

function resetPageHistory() {
    pageStack.reset()
}

async function playVideo(filepath) {
    createVideoElement({
        name: "videoPlayer",
        src: new URL(`/api/videos/stream/${filepath}`, BASE_URL)
    });
}

async function createVideoElement({ name, src }) {
    const container = document.getElementById(name);

    container.innerHTML = `
        <video controls preload="metadata">
            <source src="${src}" type="video/mp4">
            Your browser does not support the video tag.
        </video>
    `
    const videoElement = container.querySelector("video");
    videoElement.addEventListener("loadedmetadata", () => {
        const actualWidth = videoElement.videoWidth;
        const actualHeight = videoElement.videoHeight;
        console.log(`Original Resolution: ${actualWidth}x${actualHeight}`);

        videoElement.style.width = "auto";
        videoElement.style.maxWidth = "100%";
    })
    // adjustLayout(videoElement)
}

// function adjustLayout(videoElement) {
//     if (!videoElement) {
//         console.log(`videoElement: ${videoElement} not found`)
//         return;
//     }
//     videoElement.style.width = "100%";
//     videoElement.style.minWidth = "640px";
// }

// window.addEventListener("resize", () => {
//     const videoElement = document.querySelector("#videoPlayer video");
//     adjustLayout(videoElement);
// });

async function handleBackButton() {
    const previousPage = getPreviousPage()
    if (previousPage !== null) {
        await listVideos(previousPage)
    } else {
        console.log("pageStack is empty!")
    }
}

async function listVideos(path = "") {
    const response = await fetch(new URL(`/api/videos/${path}`, BASE_URL));
    const files = await response.json()
    const browser = document.getElementById("browser")
    browser.innerHTML = "";

    for (const file of files) {
        const item = document.createElement("div")
        item.textContent = file.name

        item.onclick = async () => {
            setPreviousPage(path)
            if (file.isdir) {
                await listVideos(file.filepath)
            } else {
                await playVideo(file.filepath)
            }
        }
        browser.appendChild(item)
    }
}

async function render() {
    document.getElementById("backButton").addEventListener("click", handleBackButton)
    await listVideos();
}

render()