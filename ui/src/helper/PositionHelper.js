export function getMagnitude(x, y) {
    return Math.sqrt(x * x + y * y);
}

export function normalize(x, y) {
    let magnitude = getMagnitude(x, y);
    if (magnitude > 0) {
        magnitude = magnitude / 5;
        return {x: x / magnitude, y: y / magnitude};
    } else {
        return {x: x, y: y}
    }
}