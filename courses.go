package itswizard_m_berlinPreVersion

import (
	"encoding/json"
	itswizard_basic "github.com/itslearninggermany/itswizard_m_basic"
	imses "github.com/itslearninggermany/itswizard_m_imses"
	"log"
	"time"
)

type LusdCourse struct {
	KursBezeichnung   string `json:"kursBezeichnung"`
	KursDtGeaendertAm string `json:"kursDtGeaendertAm"`
	KursFach          string `json:"kursFach"`
	KursJahrgang      string `json:"kursJahrgang"`
	KursSchulform     string `json:"kursSchulform"`
	KursUID           string `json:"kursUID"`
	SchuleUID         string `json:"schuleUID"`
}

func courseParse(out []byte, err error) (lusdCourses []LusdCourse, err2 error) {
	if err != nil {
		return nil, err
	}
	err2 = json.Unmarshal(out, &lusdCourses)
	return lusdCourses, err2
}

/*
Returns all Courses from lusd
*/
func (p *BlusdConnection) GetAllCourses() (lusdCourses []LusdCourse, out []byte, err error) {
	out, err = p.callAPI(time.Now(), "", course, false)
	if err != nil {
		return
	}
	lusdCourses, err = courseParse(out, err)
	return
}

/*
Returns all new Courses created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllCoursesCreatedScinse(minutes time.Duration) (lusdCourses []LusdCourse, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return courseParse(p.callAPI(newT, "", course, true))
}

/*
Returns all Courses from a school
*/
func (p *BlusdConnection) GetAllCoursesFromSchool(sid string) (lusdCourses []LusdCourse, err error) {
	return courseParse(p.callAPI(time.Now(), sid, course, false))
}

/*
Returns all new Courses from a school created in the last 5 minutes
*/
func (p *BlusdConnection) GetAllCoursesFromSchoolScince(sid string, minutes time.Duration) (lusdCourses []LusdCourse, err error) {
	newT := time.Now().Add(-time.Minute * minutes)
	log.Println(newT.Format(timeLayout))
	return courseParse(p.callAPI(newT, sid, course, true))
}

/*
Returns all Courses from a school created since a given time
*/
func (p *BlusdConnection) GetAllCoursesFromSchoolScinceAGivenTime(sid string, t time.Time) (lusdCourses []LusdCourse, err error) {
	return courseParse(p.callAPI(t, sid, course, true))
}

/*
Returns all Courses created since a given time
*/
func (p *BlusdConnection) GetAllCoursesScinceAGivenTime(t time.Time) (lusdCourses []LusdCourse, err error) {
	return courseParse(p.callAPI(t, "", course, true))
}

/*
func (p *blusdConnection) GetAllCoursesScinceAGivenTime (t time.Time) (lusdCourses []LusdCourse, err error){
	return courseParse(p.callAPI(t,"",course,true))
}
*/

/*
....
*/
func (p *BlusdConnection) ImportAndUpdateAllCoursesToItslearning(itsl *imses.Request) (logContent []blusdLog, err error) {
	lusdCourses, _, err := p.GetAllCourses()
	if err != nil {
		return logContent, err
	}
	for _, course := range lusdCourses {
		/*
			Todo: Schauen, ob es die Schule schon gibt.
		*/
		tmp := course.KursBezeichnung
		res := ""
		var gro []string
		for i, r := range tmp {
			res = res + string(r)
			if i > 0 && (i+1)%64 == 0 {
				gro = append(gro, res)
				res = ""
			}
		}
		courseName := res
		if len(gro) > 0 {
			courseName = gro[0]
			log.Println("gro:", gro[0])
		}

		resp, err := itsl.CreateCourse(itswizard_basic.DbGroup15{
			SyncID:        course.KursUID,
			Name:          courseName,
			ParentGroupID: course.SchuleUID + "-COURSES",
		})

		log.Println(courseName)
		if err != nil {
			logContent = append(logContent, blusdLog{
				Success: true,
				Date:    time.Now(),
				Message: "Course with id " + course.SchuleUID + " and name " + courseName + " was not (!) created in itslearning. " + resp,
			})
		} else {
			logContent = append(logContent, blusdLog{
				Success: true,
				Date:    time.Now(),
				Message: "School with id " + course.SchuleUID + " and name " + courseName + " was created in itslearning.",
			})
		}
	}
	return logContent, nil
}
