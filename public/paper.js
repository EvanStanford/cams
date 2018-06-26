/**
 * Created by valentid on 24/06/2018.
 */
var path = new Path({
    strokeColor: 'black',
    closed: true,
    fullySelected: true
});
var leftCam = new Path({
    strokeColor: 'blue',
    closed: true,
    fullySelected: true
});
var rightCam = new Path({
    strokeColor: 'red',
    closed: true,
    fullySelected: true
});

globals.path = path;
globals.leftCam = leftCam;
globals.rightCam = rightCam;

function onMouseDown(event) {
    path.add(event.point);
}