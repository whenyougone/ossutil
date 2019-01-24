package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	. "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) probeDownloadUrl(c *C, downloadFile string, repeatDown bool) {
	// create a bucket
	bucketName := bucketNamePrefix + randLowStr(12)
	_, err := s.putBucketWithACL(bucketName, "public-read")
	c.Assert(err, IsNil)

	// put a object
	objectContext := randLowStr(10240)
	uploadFileName := "ossutil_test." + randLowStr(12)
	s.createFile(uploadFileName, objectContext, c)
	object := randLowStr(12)
	s.putObject(bucketName, object, uploadFileName, c)

	// get endpoint
	tripEndpoint := ""
	pSlice := strings.Split(endpoint, "//")
	if len(pSlice) == 1 {
		tripEndpoint = pSlice[0]
	} else {
		tripEndpoint = pSlice[1]
	}

	// get object url
	// http://test-probe.oss-cn-shenzhen.aliyuncs.com/newempty1.jpg
	downUrl := "http://" + bucketName + "." + tripEndpoint + "/" + object

	var pbArgs []string
	if downloadFile == "" {
		pbArgs = []string{}
	} else {
		pbArgs = []string{downloadFile}
	}

	download := true
	options := OptionMapType{
		OptionDownload: &download,
		OptionUrl:      &downUrl,
	}

	// run command
	_, err = cm.RunCommand("probe", pbArgs, options)
	c.Assert(err, IsNil)
	c.Assert((probeCommand.pbOption.dlFilePath == ""), Equals, false)

	// check download file
	fmt.Printf("read text from %s\n", probeCommand.pbOption.dlFilePath)
	fileBody, err := ioutil.ReadFile(probeCommand.pbOption.dlFilePath)
	c.Assert(err, IsNil)
	c.Assert(objectContext, Equals, string(fileBody))

	// repeate download
	if repeatDown {
		probeCommand.pbOption.netAddr = "www.aliyun.com"
		_, err = cm.RunCommand("probe", pbArgs, options)
		c.Assert(err, IsNil)
		c.Assert((probeCommand.pbOption.dlFilePath == ""), Equals, false)
	}

	// remove bucket、file
	os.Remove(probeCommand.pbOption.dlFilePath)
	os.Remove(probeCommand.pbOption.logName)
	os.Remove(uploadFileName)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestProbeDownloadUrl(c *C) {
	// empty file name
	s.probeDownloadUrl(c, "", true)

	// file name
	s.probeDownloadUrl(c, randLowStr(12), false)

	// dir name
	dirName := "." + string(os.PathSeparator) + randLowStr(12) + string(os.PathSeparator)
	s.probeDownloadUrl(c, dirName, false)
	os.Remove(dirName)

	// exist dir name
	dirName = "." + string(os.PathSeparator) + randLowStr(12)
	err := os.MkdirAll(dirName, 0755)
	c.Assert(err, IsNil)
	s.probeDownloadUrl(c, dirName, false)
	os.Remove(dirName)
}

func (s *OssutilCommandSuite) probeDownloadWithParameter(c *C, object string, downloadFile string, repeatDown bool) {
	// create a bucket
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	var objectContext string
	var uploadFileName string
	if object != "" {
		objectContext = randLowStr(10240)
		uploadFileName = "ossutil_test." + randLowStr(12)
		s.createFile(uploadFileName, objectContext, c)
		s.putObject(bucketName, object, uploadFileName, c)
	}

	var pbArgs []string
	if downloadFile == "" {
		pbArgs = []string{}
	} else {
		pbArgs = []string{downloadFile}
	}

	download := true
	options := OptionMapType{
		OptionConfigFile: &configFile,
		OptionDownload:   &download,
		OptionBucketName: &bucketName,
	}

	if object != "" {
		options[OptionObject] = &object
	}

	tempStr := ""
	options[OptionAccessKeyID] = &tempStr
	options[OptionAccessKeySecret] = &tempStr
	options[OptionEndpoint] = &tempStr

	// run command
	_, err := cm.RunCommand("probe", pbArgs, options)
	c.Assert(err, IsNil)
	c.Assert((probeCommand.pbOption.dlFilePath == ""), Equals, false)

	// check download file
	fmt.Printf("read text from %s\n", probeCommand.pbOption.dlFilePath)
	fileBody, err := ioutil.ReadFile(probeCommand.pbOption.dlFilePath)
	c.Assert(err, IsNil)
	c.Assert(len(fileBody) > 0, Equals, true)

	if object != "" {
		c.Assert(objectContext, Equals, string(fileBody))
	}

	if repeatDown {
		probeCommand.pbOption.netAddr = "www.aliyun.com"
		_, err = cm.RunCommand("probe", pbArgs, options)
		c.Assert(err, IsNil)
		c.Assert((probeCommand.pbOption.dlFilePath == ""), Equals, false)
	}

	// remove bucket、file
	os.Remove(probeCommand.pbOption.dlFilePath)
	os.Remove(probeCommand.pbOption.logName)

	if uploadFileName != "" {
		os.Remove(uploadFileName)
	}

	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestProbeDownloadWithParameter(c *C) {
	s.probeDownloadWithParameter(c, "", "", false)
	s.probeDownloadWithParameter(c, "", randLowStr(12), false)
	s.probeDownloadWithParameter(c, randLowStr(12), "", false)
	s.probeDownloadWithParameter(c, randLowStr(12), randLowStr(12), true)
}

func (s *OssutilCommandSuite) probeUploadObject(c *C, object string, uploadFileName string, upMode string) {
	// create a bucket
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	if uploadFileName != "" {
		objectContext := randLowStr(10240)
		s.createFile(uploadFileName, objectContext, c)
	}

	var pbArgs []string
	if uploadFileName == "" {
		pbArgs = []string{}
	} else {
		pbArgs = []string{uploadFileName}
	}

	upload := true
	options := OptionMapType{
		OptionConfigFile: &configFile,
		OptionUpload:     &upload,
		OptionBucketName: &bucketName,
	}

	tempStr := ""
	options[OptionAccessKeyID] = &tempStr
	options[OptionAccessKeySecret] = &tempStr
	options[OptionEndpoint] = &tempStr

	if object != "" {
		options[OptionObject] = &object
	}

	if upMode != "" {
		options[OptionUpMode] = &upMode
	}

	// run command
	_, err := cm.RunCommand("probe", pbArgs, options)
	c.Assert(err, IsNil)

	// remove bucket、file
	os.Remove(probeCommand.pbOption.logName)

	if uploadFileName != "" {
		os.Remove(uploadFileName)
	}
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestProbeUploadObject(c *C) {
	s.probeUploadObject(c, "", "", "")
	s.probeUploadObject(c, "", randLowStr(12), "")
	s.probeUploadObject(c, "", randLowStr(12), "append")
	s.probeUploadObject(c, "", randLowStr(12), "multipart")
	s.probeUploadObject(c, randLowStr(12), randLowStr(12), "")
}

func (s *OssutilCommandSuite) TestProbeCaseError(c *C) {
	bucketName := randLowStr(12)
	testFileName := randLowStr(12)
	upload := true
	pbArgs := []string{}

	{
		options := OptionMapType{
			OptionConfigFile: &testFileName,
			OptionUpload:     &upload,
			OptionBucketName: &bucketName,
		}

		// run command
		_, err := cm.RunCommand("probe", pbArgs, options)
		c.Assert(err, NotNil)
	}

	{
		options := OptionMapType{
			OptionConfigFile: &configFile,
			OptionBucketName: &bucketName,
		}

		// run command
		_, err := cm.RunCommand("probe", pbArgs, options)
		c.Assert(err, NotNil)
	}

	{
		options := OptionMapType{
			OptionConfigFile: &configFile,
			OptionBucketName: &bucketName,
			OptionDownload:   &upload,
			OptionUpload:     &upload,
		}

		// run command
		_, err := cm.RunCommand("probe", pbArgs, options)
		c.Assert(err, NotNil)
	}
}

func (s *OssutilCommandSuite) TestProbeDownObjectError(c *C) {
	// create a bucket
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	object := randLowStr(12)
	pbArgs := []string{}

	download := true
	options := OptionMapType{
		OptionConfigFile: &configFile,
		OptionDownload:   &download,
		OptionBucketName: &bucketName,
		OptionObject:     &object,
	}

	// run command
	_, err := cm.RunCommand("probe", pbArgs, options)
	c.Assert(err, NotNil)

	// remove bucket、file
	os.Remove(probeCommand.pbOption.logName)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestProbeDownUrlError(c *C) {
	// create a bucket
	bucketName := bucketNamePrefix + randLowStr(12)
	_, err := s.putBucketWithACL(bucketName, "public-read")
	c.Assert(err, IsNil)

	// object name
	object := randLowStr(12)

	// get endpoint
	tripEndpoint := ""
	pSlice := strings.Split(endpoint, "//")
	if len(pSlice) == 1 {
		tripEndpoint = pSlice[0]
	} else {
		tripEndpoint = pSlice[1]
	}

	// get object url
	// http://test-probe.oss-cn-shenzhen.aliyuncs.com/newempty1.jpg
	downUrl := "http://" + bucketName + "." + tripEndpoint + "/" + object

	pbArgs := []string{}
	download := true
	options := OptionMapType{
		OptionDownload: &download,
		OptionUrl:      &downUrl,
	}

	// run command
	_, err = cm.RunCommand("probe", pbArgs, options)
	c.Assert(err, NotNil)

	// remove bucket、file
	os.Remove(probeCommand.pbOption.logName)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestProbeUploadObjectBucketInvalid(c *C) {
	bucketName := bucketNamePrefix + randLowStr(12)

	uploadFileName := randLowStr(12)
	objectContext := randLowStr(10240)
	s.createFile(uploadFileName, objectContext, c)

	pbArgs := []string{uploadFileName}
	upload := true
	options := OptionMapType{
		OptionConfigFile: &configFile,
		OptionUpload:     &upload,
		OptionBucketName: &bucketName,
	}

	tempStr := ""
	options[OptionAccessKeyID] = &tempStr
	options[OptionAccessKeySecret] = &tempStr
	options[OptionEndpoint] = &tempStr

	// run command
	_, err := cm.RunCommand("probe", pbArgs, options)
	c.Assert(err, NotNil)

	os.Remove(probeCommand.pbOption.logName)
	os.Remove(uploadFileName)

}

func (s *OssutilCommandSuite) TestProbeDownUrlInvalidParameter(c *C) {
	{
		probeCommand.pbOption.fromUrl = "http://test.jpg"
		err := probeCommand.downloadWithHttpUrl()
		c.Assert(err, NotNil)
	}

	{
		probeCommand.pbOption.fromUrl = "http://test-bucket/test.jpg"
		probeCommand.command.args = append(probeCommand.command.args, "oss://test.jpg")
		err := probeCommand.downloadWithHttpUrl()
		c.Assert(err, NotNil)
	}

	{
		probeCommand.pbOption.fromUrl = "http://test-bucket/test.jpg"
		probeCommand.command.args = []string{}

		err := probeCommand.downloadWithHttpUrl()
		c.Assert(err, NotNil)
	}

	{
		probeCommand.pbOption.fromUrl = "http://test-bucket/test.jpg"
		test_dir := string(os.PathSeparator) + "test-probe-" + randLowStr(5) + string(os.PathSeparator)
		probeCommand.command.args = []string{test_dir}
		err := probeCommand.downloadWithHttpUrl()
		c.Assert(err, NotNil)
	}

	probeCommand.pbOption.fromUrl = ""
	probeCommand.command.args = []string{}
}

func (s *OssutilCommandSuite) TestProbeUploadInvalidParameter(c *C) {
	{
		probeCommand.pbOption.upMode = "unkown"
		err := probeCommand.probeUpload()
		c.Assert(err, NotNil)
	}

	{
		probeCommand.pbOption.upMode = ""
		probeCommand.pbOption.bucketName = ""
		err := probeCommand.probeUpload()
		c.Assert(err, NotNil)
	}

	{
		probeCommand.pbOption.upMode = ""
		probeCommand.pbOption.bucketName = ""
		err := probeCommand.probeUpload()
		c.Assert(err, NotNil)
	}

	{
		probeCommand.pbOption.upMode = ""
		probeCommand.pbOption.bucketName = randLowStr(12)
		err := probeCommand.probeUpload()
		c.Assert(err, NotNil)
	}

	{
		probeCommand.pbOption.upMode = ""
		probeCommand.pbOption.bucketName = randLowStr(12)

		options := OptionMapType{
			OptionConfigFile: &configFile,
		}

		fileName := randLowStr(12)
		object := randLowStr(12)

		pbArgs := []string{fileName}
		probeCommand.Init(pbArgs, options)

		err := probeCommand.probeUploadFileAppend(fileName, object)
		c.Assert(err, NotNil)

		err = probeCommand.probeUploadFileMultiPart(fileName, object)
		c.Assert(err, NotNil)
	}

	{
		probeCommand.pbOption.upMode = ""
		probeCommand.pbOption.bucketName = randLowStr(12)

		options := OptionMapType{
			OptionConfigFile: &configFile,
		}

		fileName := randLowStr(12)
		object := randLowStr(12)

		pbArgs := []string{fileName}
		probeCommand.Init(pbArgs, options)

		err := probeCommand.probeUploadFileAppend(fileName, object)
		c.Assert(err, NotNil)

		// delete endpoint
		delete(probeCommand.command.options, OptionAccessKeySecret)

		err = probeCommand.probeUploadFileAppend(fileName, object)
		c.Assert(err, NotNil)

		err = probeCommand.probeUploadFileMultiPart(fileName, object)
		c.Assert(err, NotNil)

		err = probeCommand.probeUploadFileNormal(fileName, object)
		c.Assert(err, NotNil)
	}

	{
		probeCommand.pbOption.upMode = ""
		probeCommand.pbOption.bucketName = randLowStr(12)

		options := OptionMapType{
			OptionConfigFile: &configFile,
		}

		fileName := randLowStr(12)
		pbArgs := []string{fileName}
		probeCommand.Init(pbArgs, options)

		delete(probeCommand.command.options, OptionEndpoint)

		err := probeCommand.probeUpload()
		c.Assert(err, NotNil)
	}

	{
		probeCommand.pbOption.upMode = ""
		bucketName := randLowStr(12)
		probeCommand.pbOption.bucketName = bucketName

		options := OptionMapType{
			OptionConfigFile: &configFile,
		}

		fileName := randLowStr(12)
		pbArgs := []string{fileName}
		probeCommand.Init(pbArgs, options)

		tempPoint, _ := probeCommand.command.getEndpoint(bucketName)
		if !strings.Contains(tempPoint, "http") {
			tempPoint = "http://" + tempPoint
		}
		probeCommand.command.options[OptionEndpoint] = tempPoint

		probeCommand.command.args = []string{}
		probeCommand.command.args = append(probeCommand.command.args, "oss://temp-probe")

		err := probeCommand.probeUpload()
		c.Assert(err, NotNil)
	}
}

func (s *OssutilCommandSuite) TestProbeUploadObjectRepeat(c *C) {
	// create a bucket
	bucketName := bucketNamePrefix + randLowStr(12)
	s.putBucket(bucketName, c)

	uploadFileName := randLowStr(12)
	objectContext := randLowStr(10240)
	s.createFile(uploadFileName, objectContext, c)

	object := randLowStr(12)
	s.putObject(bucketName, object, uploadFileName, c)

	pbArgs := []string{uploadFileName}
	upload := true
	options := OptionMapType{
		OptionConfigFile: &configFile,
		OptionUpload:     &upload,
		OptionBucketName: &bucketName,
		OptionObject:     &object,
	}

	tempStr := ""
	options[OptionAccessKeyID] = &tempStr
	options[OptionAccessKeySecret] = &tempStr
	options[OptionEndpoint] = &tempStr

	// run command
	_, err := cm.RunCommand("probe", pbArgs, options)
	c.Assert(err, IsNil)

	// remove bucket、file
	os.Remove(probeCommand.pbOption.logName)
	os.Remove(uploadFileName)
	s.removeBucket(bucketName, true, c)
}

func (s *OssutilCommandSuite) TestProbeDownloadInvalidParameter(c *C) {
	{
		probeCommand.pbOption.fromUrl = "http/////"
		err := probeCommand.probeDownload()
		c.Assert(err, NotNil)
	}

	{
		probeCommand.pbOption.fromUrl = ""
		probeCommand.pbOption.bucketName = ""
		err := probeCommand.probeDownload()
		c.Assert(err, NotNil)
	}

	{
		probeCommand.command.args = []string{}
		probeCommand.command.args = append(probeCommand.command.args, "oss:////")
		_, err := probeCommand.getFileNameArg()
		c.Assert(err, NotNil)
		probeCommand.command.args = []string{}
	}
}

func (s *OssutilCommandSuite) TestProbeDownloadObjectInvalidParameter(c *C) {
	{
		probeCommand.command.args = []string{}
		probeCommand.command.args = append(probeCommand.command.args, "oss://test-probe/file.txt")

		var srcURL CloudURL
		srcURL.bucket = randLowStr(10)
		srcURL.object = randLowStr(10)

		err := probeCommand.probeDownloadObject(srcURL, true)
		c.Assert(err, NotNil)
	}

	{
		probeCommand.command.args = []string{}
		var srcURL CloudURL
		srcURL.bucket = randLowStr(10)
		srcURL.object = randLowStr(10)

		err := probeCommand.probeDownloadObject(srcURL, true)
		c.Assert(err, NotNil)
	}
}
