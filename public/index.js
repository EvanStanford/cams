/**
 * Created by valentid on 24/06/2018.
 */
var globals = {};

window.onload = function() {
    convert.onclick = function(){
        let path = [];
        globals.path.segments.forEach(function(segment){
            path.push({x:segment.point.x,y:segment.point.y});
        });

        let cams = createCams(path,43.5, 32.3, 7.06);//units = mm
        drawCams(cams);

        BinaryStlWriter.save(toSTL(cams.cams.left, cams.centers.left), "leftCam.stl");
        BinaryStlWriter.save(toSTL(cams.cams.right, cams.centers.right), "rightCam.stl");
    };
};

function toSTL(cam, center){
    //triangulate a pizza out of the cam polygon
    let geometry = {faces: [], vertices: cam.map(p=>{return {x:p.x,y:p.y,z:0};})};
    geometry.vertices.push({x:center.x, y:center.y, z:0});//the center sits in vertices[cam.length]
    for(let i=0; i<cam.length; i++){
        geometry.faces.push({a:i, b:(i+1)%cam.length, c:cam.length, normal:{x:0,y:0,z:1}});
    }
    return geometry;
}

function drawCams(cams){
    let leftCamPaper = new paper.PaperScope();
    let rightCamPaper = new paper.PaperScope();
    leftCamPaper.setup(document.getElementById('leftCamCanvas'));
    rightCamPaper.setup(document.getElementById('rightCamCanvas'));
    drawCamIntoPaper(cams.cams.left, cams.centers.left, leftCamPaper);
    drawCamIntoPaper(cams.cams.right, cams.centers.right, rightCamPaper);
}

function drawCamIntoPaper(path, center, paper){
    paper.activate();
    let normalizedPath = path.slice();
    normalizedPath.push(center);
    normalizedPath = normalize(normalizedPath,200);//local deep copy of path
    let normalizedCenter = normalizedPath.pop();
    let _path = new paper.Path({strokeColor:'black'});
    _path.moveTo(new paper.Point(normalizedPath[0].x+200,normalizedPath[0].y+200));
    for(let i=1; i<=normalizedPath.length; ++i){
        _path.lineTo(new paper.Point(normalizedPath[i%normalizedPath.length].x+200,normalizedPath[i%normalizedPath.length].y+200));
    }
    new paper.Path.Circle(new paper.Point(normalizedCenter.x+200,normalizedCenter.y+200),5).fillColor = 'black';
    paper.view.draw();
}


/**
 * 1. flips y axis s.t (up +, down -)
 * 2. centers the path around (0,0)
 * 3. scales s.t max point distance from (0,0) is == MAX_DIST
 *
 * @param path
 * @returns modified path
 */
function normalize(path, MAX_DIST){
    //flip and get center
    let centerX=0,centerY=0;
    path.forEach(function(p){
        p.y=-p.y;
        centerX+=p.x;
        centerY+=p.y;
    });
    centerX/=path.length;
    centerY/=path.length;

    //move center to (0,0) and calc maxDist
    let maxDist = 0;
    path.forEach(function(p){
        let pDist = p.x*p.x + p.y*p.y;
        if(maxDist<pDist) maxDist = pDist;
        p.x-=centerX;
        p.y-=centerY;
    });
    maxDist = Math.sqrt(maxDist);

    //normalize distance
    let factor = MAX_DIST/maxDist;
    path.forEach(function(p){
        p.x = p.x*factor;
        p.y = p.y*factor;
    });
    return path;
}

function interpolate(path, n){
    let interpolated = new Array(path.length * n);
    for(let i=0; i<path.length; ++i){
        let p1 = path[i];
        let p2 = path[(i+1)%path.length];
        let dx = (p2.x-p1.x)/n;
        let dy = (p2.y-p1.y)/n;
        for(let d=0; d<n; ++d){
            interpolated[i*n+d] = {x:p1.x+(dx*d), y:p1.y+(dy*d)};
        }
    }
    return interpolated;
}

function createCams(path, camBoreSpacing, maxCamRadius, penRadius){
    path = normalize(path, 10);//here we define that max pen movement is 10mm from its (0,0) position
    path = interpolate(path,15);

    let centers = getCamCenters(path, penRadius, camBoreSpacing, maxCamRadius);
    let cams = getCams(path, centers, penRadius);

    return {cams:cams, centers:centers};
}

/**
 * binary search y s.t max radius
 * @param path
 * @param penRadius - Laser pen radius
 * @param camBoreSpacing - distance between cam axels, center to center
 * @param maxCamRadius - max radius of cam. If this was any larger the cams would hit each other
 */
function getCamCenters(path, penRadius, camBoreSpacing, maxCamRadius){
    let minCamY = - maxCamRadius;
    let maxCamY = 0;
    let leftCenter, rightCenter;
    for(let i=0; i<15; ++i) {
        let testCamY = (minCamY + maxCamY) / 2;
        leftCenter = {x: (-camBoreSpacing / 2), y: testCamY};
        rightCenter = {x: (camBoreSpacing / 2), y: testCamY};
        if (maxCamRadius < getMaxRadius(path, leftCenter, rightCenter, penRadius)) {
            minCamY = testCamY;
        } else {
            maxCamY = testCamY;
        }
    }
    return {left: leftCenter, right:rightCenter};
}

function getMaxRadius(path, leftCamCenter, rightCamCenter, penRadius) {
    let maxR = 0;
    path.forEach(function (target){
        let r;
        r = getRadius(target, leftCamCenter, penRadius);
        if(r > maxR) maxR = r;
        r = getRadius(target, rightCamCenter, penRadius);
        if(r > maxR) maxR = r;
    });
    return maxR;
}

function getRadius(target, camCenter, penRadius){
    let dx = target.x - camCenter.x;
    let dy = target.y - camCenter.y;
    return Math.sqrt(dx*dx+dy*dy)-penRadius
}

function getCams(path, camCenters, penRadius){
    let left = Array(path.length);
    let right = Array(path.length);

    for(let i=0; i<path.length; ++i){
        let target = path[i];
        left[i] = convertCoord(target, camCenters.left, path.length, i, penRadius, false);
        right[i] = convertCoord(target, camCenters.right, path.length, i, penRadius, true);
    }

    return {left:left, right:right};
}

/**
 *
 * @param target
 * @param camCenter
 * @param n - total number of coordinates
 * @param i - index of this coordinate
 * @param penRadius
 * @param right - is this the right cam?
 * @returns {{x: number, y: number}}
 */
function convertCoord(target, camCenter, n, i, penRadius, right) {
    let rx = getRadius(target, camCenter, penRadius);
    let B = Math.atan((target.y-camCenter.y)/(target.x-camCenter.x));
    let Kx = 2*Math.PI*i/n;
    if(right) Kx = -(Kx+Math.PI);
    let θ = 2*Math.PI-Kx+B;
    return {x:(rx*Math.cos(θ)+camCenter.x), y:(rx*Math.sin(θ)+camCenter.y)};
}